// Copyright 2017 Google LLC
// Copyright 2023 Benno Van Waeyenberg
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//	Fork of cloud.google.com/go/firestore with the following changes:
// 	* Removed dependencies on the Firestore API
// 	* Removed dependency on Firestore Protobuf value types

package firestruct

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"time"

	ts "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"

	"github.com/bennovw/firestruct/internal/fields"

	"google.golang.org/genproto/googleapis/type/latlng"
)

var (
	typeOfByteSlice      = reflect.TypeOf([]byte{})
	typeOfGoTime         = reflect.TypeOf(time.Time{})
	typeOfLatLng         = reflect.TypeOf(latlng.LatLng{})
	typeOfUUID           = reflect.TypeOf(uuid.UUID{})
	typeOfProtoTimestamp = reflect.TypeOf((*ts.Timestamp)(nil))
)

// dataToReflectPointer uses any type of value to set p, which should be a pointer to a struct.
// An error is returned if data value types don't match p, or if your struct fields are private (capitalize your struct fields).
// You may add tags to your struct fields formatted as `firestore:"changeme"` to specify the field name to map to. If you do not specify a tag, the field name will be used.
// If the Firestore document contains a field that is not present in the struct, it will be ignored. If the struct contains a field that is not present in the Firestore document, it will be set to its zero value.
//
// Example:
//
//	type Person struct {
//		Name string
//		Age  int
//	}
//	var p Person
//	err := dataToReflectPointer(reflect.ValueOf(p).Elem(), map[string]interface{}{"Name": "John", "Age": 21})
func dataToReflectPointer(p reflect.Value, data any) error {
	typeErr := func() error {
		return fmt.Errorf("cannot use value %T to populate %s ", data, p.Type())
	}

	// A Null value sets anything nullable to nil, and has no effect
	// on anything else.
	if data == nil {
		switch p.Kind() {
		case reflect.Interface, reflect.Ptr, reflect.Map, reflect.Slice:
			p.Set(reflect.Zero(p.Type()))
		}
		return nil
	}

	// Handle special types first.
	switch p.Type() {
	case typeOfByteSlice:
		switch x := data.(type) {
		case string:
			b, err := base64.StdEncoding.DecodeString(x)
			if err != nil {
				return typeErr()
			}
			p.SetBytes(b)
			return nil

		case []byte:
			p.SetBytes(x)
			return nil

		default:
			return typeErr()
		}

	case typeOfGoTime:
		x, ok := data.(string)
		if !ok {
			return typeErr()
		}

		ts, err := time.Parse(time.RFC3339, x)
		if err != nil {
			return typeErr()
		}

		p.Set(reflect.ValueOf(ts))
		return nil

	case typeOfLatLng:
		switch x := data.(type) {
		case latlng.LatLng:
			p.Set(reflect.ValueOf(x))
			return nil

		case map[string]interface{}:
			lat, ok := x["latitude"].(float64)
			if !ok {
				return errors.New("latitude is not a float64")
			}
			lng, ok := x["longitude"].(float64)
			if !ok {
				return errors.New("longitude is not a float64")
			}
			p.Set(reflect.ValueOf(latlng.LatLng{Latitude: lat, Longitude: lng}))
			return nil

		default:
			return typeErr()
		}

	case typeOfUUID:
		x, ok := data.(string)
		if !ok {
			return typeErr()
		}
		uuid, err := uuid.Parse(x)
		if err != nil {
			return fmt.Errorf("%v is not a valid UUID: %v", data, err)

		}
		p.Set(reflect.ValueOf(uuid))
		return nil

	}

	switch p.Kind() {
	case reflect.Bool:
		x, ok := data.(bool)
		if !ok {
			return typeErr()
		}
		p.SetBool(x)

	case reflect.String:
		x, ok := data.(string)
		if !ok {
			return typeErr()
		}
		p.SetString(x)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int64
		switch x := data.(type) {
		case int:
			i = int64(x)
		case int8:
			i = int64(x)
		case int16:
			i = int64(x)
		case int32:
			i = int64(x)
		case int64:
			i = x
		case float64:
			i = int64(x)
		default:
			return typeErr()
		}

		if p.OverflowInt(i) {
			return overflowErr(p, data)
		}
		p.SetInt(i)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32:
		var u uint64
		switch x := data.(type) {
		case uint8:
			u = uint64(x)
		case uint16:
			u = uint64(x)
		case uint32:
			u = uint64(x)
		case uint64:
			u = x
		default:
			return typeErr()
		}

		if p.OverflowUint(u) {
			return overflowErr(p, data)
		}
		p.SetUint(u)

	case reflect.Float32, reflect.Float64:
		var f float64
		switch x := data.(type) {
		case float32:
			f = float64(x)
		case float64:
			f = x
		default:
			return typeErr()
		}

		if p.OverflowFloat(f) {
			return overflowErr(p, data)
		}
		p.SetFloat(f)

	case reflect.Slice:
		vals, ok := data.([]any)
		if !ok {
			return typeErr()
		}
		vlen := p.Len()
		xlen := len(vals)
		// Make a slice of the right size, avoiding allocation if possible.
		switch {
		case vlen < xlen:
			p.Set(reflect.MakeSlice(p.Type(), xlen, xlen))
		case vlen > xlen:
			p.SetLen(xlen)
		}
		return populateArray(p, vals, xlen)

	case reflect.Array:
		vals, ok := data.([]any)
		if !ok {
			return typeErr()
		}

		xlen := len(vals)
		vlen := p.Len()
		minlen := vlen
		// Set extra elements to their zero value.
		if vlen > xlen {
			z := reflect.Zero(p.Type().Elem())
			for i := xlen; i < vlen; i++ {
				p.Index(i).Set(z)
			}
			minlen = xlen
		}
		return populateArray(p, vals, minlen)

	case reflect.Map:
		x, ok := data.(map[string]any)
		if !ok {
			return typeErr()
		}

		return populateMap(p, x)

	case reflect.Ptr:
		// If the pointer is nil, set it to a zero value.
		if p.IsNil() {
			p.Set(reflect.New(p.Type().Elem()))
		}
		return dataToReflectPointer(p.Elem(), data)

	case reflect.Struct:
		x, ok := data.(map[string]any)
		if !ok {
			return typeErr()
		}
		return populateStruct(p, x)

	case reflect.Interface:
		if p.NumMethod() == 0 { // empty interface
			// If p holds a pointer, set the pointer.
			if !p.IsNil() && p.Elem().Kind() == reflect.Ptr {
				return dataToReflectPointer(p.Elem(), data)
			}
			// Otherwise, create a fresh value.
			p.Set(reflect.ValueOf(data))
			return nil
		}
		// Any other kind of interface is an error.
		fallthrough

	default:
		return fmt.Errorf("cannot set type %s", p.Type())
	}
	return nil
}

// populateArray sets the first n elements of vr, which must be a slice or
// array, to the corresponding elements of vals.
func populateArray(vr reflect.Value, vals []any, n int) error {
	for i := 0; i < n; i++ {
		if err := dataToReflectPointer(vr.Index(i), vals[i]); err != nil {
			return err
		}
	}
	return nil
}

// populateMap sets the elements of vm, which must be a map, from the
// corresponding elements of pm.
//
// Since a map value is not settable, this function always creates a new
// element for each corresponding map key. Existing values of vm are
// overwritten. This happens even if the map value is something like a pointer
// to a struct, where we could in theory populate the existing struct value
// instead of discarding it. This behavior matches encoding/json.
func populateMap(vm reflect.Value, pm map[string]any) error {
	t := vm.Type()
	if t.Key().Kind() != reflect.String {
		return errors.New("map key type is not string")
	}
	if vm.IsNil() {
		vm.Set(reflect.MakeMap(t))
	}
	et := t.Elem()
	for k, vproto := range pm {
		el := reflect.New(et).Elem()
		if err := dataToReflectPointer(el, vproto); err != nil {
			return err
		}
		vm.SetMapIndex(reflect.ValueOf(k), el)
	}
	return nil
}

// populateStruct sets the fields of vs, which must be a struct, from
// the matching elements of pm.
func populateStruct(vs reflect.Value, data map[string]any) error {
	fs, err := fieldCache.Fields(vs.Type())
	if err != nil {
		return err
	}

	type match struct {
		val any
		f   *fields.Field
	}
	// Find best field matches
	matched := make(map[string]match)
	for k, field := range data {
		f := fs.Match(k)
		if f == nil {
			continue
		}
		if _, ok := matched[f.Name]; ok {
			// If multiple case insensitive fields match, the exact match
			// should win.
			if f.Name == k {
				matched[k] = match{val: field, f: f}
			}
		} else {
			matched[f.Name] = match{val: field, f: f}
		}
	}

	// Reflect values
	for _, v := range matched {
		f := v.f
		val := v.val

		if err := dataToReflectPointer(vs.FieldByIndex(f.Index), val); err != nil {
			return fmt.Errorf("%s.%s: %w", vs.Type(), f.Name, err)
		}
	}
	return nil
}

func overflowErr(v reflect.Value, x interface{}) error {
	return fmt.Errorf("value %v overflows type %s", x, v.Type())
}

var fieldCache = fields.NewCache(parseTag, nil, isLeafType)

// isLeafType determines whether or not a type is a 'leaf type'
// and should not be recursed into, but considered one field.
func isLeafType(t reflect.Type) bool {
	return t == typeOfGoTime || t == typeOfLatLng || t == typeOfProtoTimestamp
}
