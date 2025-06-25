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

package firestruct

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"

	"google.golang.org/genproto/googleapis/type/latlng"
)

var (
	protoBytesTag     = "bytesValue"
	protoMapTag       = "mapValue"
	protoArrayTag     = "arrayValue"
	protoNullTag      = "nullValue"
	protoStringTag    = "stringValue"
	protoBoolTag      = "booleanValue"
	protoIntTag       = "integerValue"
	protoDoubleTag    = "doubleValue"
	protoTimestampTag = "timestampValue"
	protoGeoPointTag  = "geoPointValue"
	protoReferenceTag = "referenceValue"
)

// List of Firestore protojson tags without any nested data structures
// For full listing of Firestore protojson tags, see https://firebase.google.com/docs/firestore/reference/rest/v1/Value
var FirestoreFlatDataTypes = []string{
	protoStringTag,
	protoBoolTag,
	protoReferenceTag,
	protoTimestampTag,
	protoNullTag,
	protoIntTag,
	protoDoubleTag,
	protoBytesTag,
	protoGeoPointTag,
}

var FirestoreSimpleDataTypes = []string{
	protoStringTag,
	protoBoolTag,
	protoReferenceTag,
	protoTimestampTag,
	protoNullTag,
}

// UnwrapFirestoreFields unwraps a map[string]any containing one or more nested Firestore protojson encoded fields and returns a Go map[string]any without Firestore protojson tags.
func UnwrapFirestoreFields(input map[string]any) (map[string]any, error) {
	if input == nil {
		return nil, errors.New("nil map contents")
	}

	output := make(map[string]any, len(input))
	emptyMap := make(map[string]interface{})
	mapType := reflect.TypeOf(emptyMap)

	for k, val := range input {
		vType := reflect.TypeOf(val)

		if vType != mapType {
			return nil, fmt.Errorf("invalid input, expecting *map[string]any, but received %T", val)
		}

		// handle less common cases first
		if len(input) == 1 {
			if k == protoMapTag {
				// if the document only contains a single map without a title descriptor, we can return the map directly
				x, err := unwrapMap(val)
				if err != nil {
					return nil, err
				}

				return x, nil
			} else if k == protoArrayTag {
				// when a document contains an array the immediate children won't have a title descriptor, so no need to unwrap the title
				x, err := unwrapArray(val)
				if err != nil {
					return nil, err
				}

				output[k] = x
				return output, nil
			}
		}

		// usually the top level of the input map is a title descriptor, we evaluate the protojson tags in the subvalues before unwrapping our data
		for kk, vv := range val.(map[string]interface{}) {
			// Process data types that don't contain nested data first
			if kk != protoMapTag && kk != protoArrayTag {
				x, err := unwrapFlatValue(val)
				if err != nil {
					return nil, err
				}
				output[k] = x

				continue
			}

			// recursively process maps
			if kk == protoMapTag {
				x, err := unwrapMap(vv)
				if err != nil {
					return nil, err
				}

				output[k] = x
			}

			// recursively process arrays as slices
			if kk == protoArrayTag {
				x, err := unwrapArray(vv)
				if err != nil {
					return nil, err
				}

				output[k] = x
			}
		}

	}

	return output, nil
}

// unwrapFlatValue unwraps shallow Firestore data types (i.e. those without nested data structures)
func unwrapFlatValue(value any) (any, error) {
	mapValue, ok := value.(map[string]interface{})
	if !ok || len(mapValue) != 1 {
		return nil, fmt.Errorf("unwrapFlatValue error processing unsupported value: %v", value)
	}

	// Check if the value in the payload is encoded, starting with bytes
	if bytesValue, ok := mapValue[protoBytesTag]; ok {
		return unwrapBytes(bytesValue)
	}

	// Ensure int values are converted from float64 to int
	if intValue, ok := mapValue[protoIntTag]; ok {
		return unwrapInt(intValue)
	}

	// Ensure float values without decimal point are converted from int to float64
	if doubleValue, ok := mapValue[protoDoubleTag]; ok {
		return unwrapDouble(doubleValue)
	}

	// Ensure geopoint values are converted from map[string]interface{} to GeoPoint
	if geoPointValue, ok := mapValue[protoGeoPointTag]; ok {
		return unwrapGeoPoint(geoPointValue)
	}

	for _, key := range FirestoreSimpleDataTypes {
		subValue, ok := mapValue[key]
		if !ok {
			continue
		}

		return subValue, nil
	}

	return nil, fmt.Errorf("unwrapFlatValue error processing unsupported value: %v", value)
}

// unwrapMap returns the values nested within a Firestore json encoded map
func unwrapMap(value any) (map[string]any, error) {
	m, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unwrapMap error, Firestore map is expected to be a map[string]interface{} got: %T", value)
	}
	mf, ok := m["fields"]
	if !ok {
		// if the map is empty, return nil
		if len(m) == 0 {
			return nil, nil
		}

		return nil, fmt.Errorf("unwrapMap error, \"fields\" key not found")
	}
	mv, ok := mf.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unwrapMap erro, Firestore map fields are expected to be a map[string]interface{} got: %T", value)
	}

	subValues, err := UnwrapFirestoreFields(mv)
	if err != nil {
		return nil, err
	}
	return subValues, nil
}

// unwrapArray returns the array values nested within a Firestore json encoded array
func unwrapArray(array any) ([]any, error) {
	am, ok := array.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unwrapArray error, Firestore array is expected to be a map[string]interface{}")
	}

	v, ok := am["values"]
	if !ok {
		// if the array is empty, return nil
		if len(am) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("unwrapArray error, \"values\" key not found")
	}

	va, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("unwrapArray error, \"values\" does not contain an array of values")
	}

	// create new array and populate it with unwrapped subvalues
	outputArray := make([]any, len(va))

	for i, val := range va {
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("unwrapArray error, array can only contain values encoded as map[string]interface{}")
		}

		output, err := UnwrapFirestoreFields(mapVal)
		if err != nil {
			return nil, err
		}
		outputArray[i] = output
	}

	return outputArray, nil
}

// unwrapBytes decodes base64-encoded bytes from Firestore protojson
func unwrapBytes(bytesValue any) ([]byte, error) {
	if bytesValue == nil {
		return nil, fmt.Errorf("unwrapBytes error: nil bytes value")
	}

	// If the bytesValue is already a byte slice, if so, return it directly
	if bv, ok := bytesValue.([]byte); ok {
		return bv, nil
	}

	// Check if the bytesValue is a string, which is how Firestore encodes bytes in protojson format
	bv, ok := bytesValue.(string)
	if ok {
		decoded, err := base64.StdEncoding.DecodeString(bv)
		if err != nil {
			return nil, fmt.Errorf("unwrapBytes error decoding base64 bytes value: %v", err)
		}
		return decoded, nil
	}

	return nil, fmt.Errorf("unwrapBytes error processing bytes value: %v", bytesValue)
}

// unwrapInt converts integer values from float64
func unwrapInt(intValue any) (int, error) {
	if iv, ok := intValue.(float64); ok {
		return int(iv), nil
	}

	if iv, ok := intValue.(int); ok {
		return iv, nil
	}

	return 0, fmt.Errorf("unwrapInt error processing int value: %v", intValue)
}

// unwrapDouble converts double values from int to float64 if they are not already float64
func unwrapDouble(doubleValue any) (float64, error) {
	if dv, ok := doubleValue.(float64); ok {
		return dv, nil
	}

	if dv, ok := doubleValue.(int); ok {
		return float64(dv), nil
	}

	return 0, fmt.Errorf("unwrapDouble error processing double value: %v", doubleValue)
}

// unwrapGeoPoint converts geopoint values from map[string]interface{} to latlng.LatLng
func unwrapGeoPoint(geoPointValue any) (latlng.LatLng, error) {
	if geoPointValue == nil {
		return latlng.LatLng{}, nil
	}

	// check if the geoPointValue is already a latlng.LatLng type, if so, return it directly
	if gp, ok := geoPointValue.(latlng.LatLng); ok {
		return gp, nil
	}

	gp, ok := geoPointValue.(map[string]interface{})
	if !ok {
		return latlng.LatLng{}, fmt.Errorf("unwrapGeoPoint error processing geoPoint value: %v", geoPointValue)
	}

	lat, ok := gp["latitude"].(float64)
	if !ok {
		return latlng.LatLng{}, fmt.Errorf("unwrapGeoPoint error processing geoPoint latitude value: %v", gp["latitude"])
	}

	lng, ok := gp["longitude"].(float64)
	if !ok {
		return latlng.LatLng{}, fmt.Errorf("unwrapGeoPoint error processing geoPoint longitude value: %v", gp["longitude"])
	}

	return latlng.LatLng{Latitude: lat, Longitude: lng}, nil
}
