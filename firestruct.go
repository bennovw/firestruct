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
// 	* Added support for unwrapping Firestore JSON encoded data types into native Go types

package firestruct

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// A Google Cloud Event fired when a Firestore document is created, updated or deleted.
// FirestoreCloudEvent is the payload of a Firestore event, it contains the current and old version of the Firestore document triggering the event.
type FirestoreCloudEvent struct {
	OldValue   FirestoreDocument `firestore:"oldValue,omitempty" json:"oldValue,omitempty"`
	Value      FirestoreDocument `firestore:"value" json:"value"`
	UpdateMask struct {
		FieldPaths []string `firestore:"fieldPaths,omitempty" json:"fieldPaths,omitempty"`
	} `firestore:"updateMask,omitempty" json:"updateMask,omitempty"`
}

// Document is an alias for Value, it returns the current version of the Firestore document triggering the event.
func (e *FirestoreCloudEvent) Document() *FirestoreDocument {
	return &e.Value
}

// DataTo uses the current version of the Firestore document to populate p, which should be a pointer to a struct or a pointer to a map[string]interface{}.
// You may add tags to your struct fields formatted as `firestore:"changeme"` to specify the Firestore field name to use. If you do not specify a tag, the field name will be used.
// If the Firestore document contains a field that is not present in the struct, it will be ignored. If the struct contains a field that is not present in the Firestore document, it will be set to its zero value.
func (e *FirestoreCloudEvent) DataTo(p interface{}) error {
	return e.Value.DataTo(p)
}

// ToMap returns the current version of the Firestore document as an unwrapped map[string]interface{} without any nested protojson type descriptor tags.
func (e *FirestoreCloudEvent) ToMap() (map[string]any, error) {
	m, err := e.Value.ToMap()
	return m, err
}

// A Firestore document.
// Fields contains Firestore JSON encoded data types, see https://Firestore.google.com/docs/firestore/reference/rest/v1/Value
type FirestoreDocument struct {
	Name       string         `firestore:"name,omitempty" json:"name,omitempty"`
	Fields     map[string]any `firestore:"fields,omitempty" json:"fields,omitempty"`
	CreateTime time.Time      `firestore:"createTime,serverTimestamp,omitempty" json:"createTime,omitempty"`
	UpdateTime time.Time      `firestore:"updateTime,serverTimestamp,omitempty" json:"updateTime,omitempty"`
}

// DataTo uses the document's fields to populate p, which can be a pointer to a
// map[string]interface{} or a pointer to a struct.
// You may add tags to your struct fields formatted as `firestore:"changeme"` to specify the Firestore field name to use. If you do not specify a tag, the field name will be used.
// If the Firestore document contains a field that is not present in the struct, it will be ignored. If the struct contains a field that is not present in the Firestore document, it will be set to its zero value.
//
// Firestore field values are converted to Go values as follows:
//   - Null converts to nil.
//   - Bool converts to bool.
//   - String converts to string.
//   - Integer converts int64. When setting a struct field, any signed or unsigned
//     integer type is permitted except uint, uint64 or uintptr. Overflow is detected
//     and results in an error.
//   - Double converts to float64. When setting a struct field, float32 is permitted.
//     Overflow is detected and results in an error.
//   - Bytes is converted to []byte.
//   - Timestamp converts to time.Time.
//   - GeoPoint converts to *latlng.LatLng, where latlng is the package
//     "google.golang.org/genproto/googleapis/type/latlng".
//   - Arrays convert to []interface{}. When setting a struct field, the field
//     may be a slice or array of any type and is populated recursively.
//     Slices are resized to the incoming value's size, while arrays that are too
//     long have excess elements filled with zero values. If the array is too short,
//     excess incoming values will be dropped.
//   - Maps convert to map[string]interface{}. When setting a struct field,
//     maps of key type string and any value type are permitted, and are populated
//     recursively.
//   - WARNING: Firestore document references are NOT handled.
//
// Field names given by struct field tags are observed, as described in
// DocumentRef.Create.
//
// Only the fields actually present in the document are used to populate p. Other fields
// of p are left unchanged.

func (d *FirestoreDocument) DataTo(p interface{}) error {
	// Remove Firestore protojson field tags from the document's fields.
	flatDoc, err := d.ToMap()
	if err != nil {
		return fmt.Errorf("firestruct: error converting Firestore document to map %v", err)
	}

	return DataTo(p, flatDoc)
}

// ToMap converts a Firestore document to a native Go map[string]interface{} without protojson tags
func (e *FirestoreDocument) ToMap() (map[string]any, error) {
	if e == nil {
		return nil, errors.New("firestruct: nil document contents")
	}

	fields, err := UnwrapFirestoreFields(e.Fields)
	if err != nil {
		return nil, err
	}
	return fields, nil
}

// DataTo uses the input data to populate p, which can be a pointer to a struct or a pointer to a map[string]interface{}.
// You may add tags to your struct fields formatted as `firestore:"changeme"` to specify the Firestore field name to use. If you do not specify a tag, the field name will be used.
// If the input data contains a field that is not present in the struct, it will be ignored. If the struct contains a field that is not present in the input data, it will be set to its zero value.
func DataTo(pointer interface{}, data any) error {
	pv := reflect.ValueOf(pointer)
	if pv.Kind() != reflect.Ptr || pv.IsNil() {
		return errors.New("firestruct: target is nil or not a pointer to a struct or map")
	}

	// If p is a pointer to a map, populate it directly.
	_, ok := pointer.(map[string]any)
	if ok {
		pv.Elem().Set(reflect.ValueOf(data))
		return nil
	}

	// Otherwise, p is a pointer to a struct, so populate it recursively.
	return dataToReflectPointer(pv.Elem(), data)
}