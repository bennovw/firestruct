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
	"errors"
	"fmt"
	"reflect"
)

// List of Firestore protojson tags without any nested data structures
// For full listing of Firestore protojson tags, see https://firebase.google.com/docs/firestore/reference/rest/v1/Value
var FirestoreFlatDataTypes = []string{
	"stringValue",
	"booleanValue",
	"integerValue",
	"doubleValue",
	"timestampValue",
	"nullValue",
	"bytesValue",
	"referenceValue",
	"geoPointValue",
}

// UnwrapFirestoreFields unwraps a map[string]any containing one or more nested Firestore protojson encoded fields and returns a Go map[string]any without Firestore protojson tags.
func UnwrapFirestoreFields(input map[string]any) (map[string]any, error) {
	if input == nil {
		return nil, errors.New("firestruct: nil map contents")
	}

	output := make(map[string]any, len(input))
	emptyMap := make(map[string]interface{})
	mapType := reflect.TypeOf(emptyMap)

	for k, val := range input {
		vType := reflect.TypeOf(val)

		if vType != mapType {
			return nil, fmt.Errorf("firestruct: invalid input, expecting *map[string]any, but received %T", val)
		}

		// handle less common cases first
		if len(input) == 1 {
			if k == "mapValue" {
				// if the document only contains a single map without a title descriptor, we can return the map directly
				x, err := unwrapMap(val)
				if err != nil {
					return nil, err
				}

				return x, nil
			} else if k == "arrayValue" {
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
			if kk != "mapValue" && kk != "arrayValue" {
				x, err := unwrapFlatValue(val)
				if err != nil {
					return nil, err
				}
				output[k] = x

				continue
			}

			// recursively process maps
			if kk == "mapValue" {
				x, err := unwrapMap(vv)
				if err != nil {
					return nil, err
				}

				output[k] = x
			}

			// recursively process arrays as slices
			if kk == "arrayValue" {
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
		return nil, fmt.Errorf("firestruct: unwrapFlatValue error processing unsupported value: %v", value)
	}

	for _, subkey := range FirestoreFlatDataTypes {
		subValue, ok := mapValue[subkey]
		if !ok {
			continue
		}

		return subValue, nil
	}

	return nil, fmt.Errorf("firestruct: unwrapFlatValue error processing unsupported value: %v", value)
}

// unwrapMap returns the values nested within a Firestore json encoded map
func unwrapMap(value any) (map[string]any, error) {
	m, ok := value.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("firestruct: unwrapMap error, Firestore map is expected to be a map[string]interface{} got: %T", value)
	}
	mf, ok := m["fields"]
	if !ok {
		// if the map is empty, return nil
		if len(m) == 0 {
			return nil, nil
		}

		return nil, fmt.Errorf("firestruct: unwrapMap error, \"fields\" key not found")
	}
	mv, ok := mf.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("firestruct: unwrapMap erro, Firestore map fields are expected to be a map[string]interface{} got: %T", value)
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
		return nil, fmt.Errorf("firestruct: unwrapArray error, Firestore array is expected to be a map[string]interface{}")
	}

	v, ok := am["values"]
	if !ok {
		// if the array is empty, return nil
		if len(am) == 0 {
			return nil, nil
		}
		return nil, fmt.Errorf("firestruct: unwrapArray error, \"values\" key not found")
	}

	va, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("firestruct: unwrapArray error, \"values\" does not contain an array of values")
	}

	// create new array and populate it with unwrapped subvalues
	outputArray := make([]any, len(va))

	for i, val := range va {
		mapVal, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("firestruct: unwrapArray error, array can only contain values encoded as map[string]interface{}")
		}

		output, err := UnwrapFirestoreFields(mapVal)
		if err != nil {
			return nil, err
		}
		outputArray[i] = output
	}

	return outputArray, nil
}
