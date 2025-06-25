// Copyright 2017 Google LLC
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

package testutil

import (
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/fatih/structs"
	"github.com/golang/protobuf/proto"
	"github.com/google/go-cmp/cmp"
)

var (
	alwaysEqual = cmp.Comparer(func(_, _ interface{}) bool { return true })

	defaultCmpOptions = []cmp.Option{
		// Use proto.Equal for protobufs
		cmp.Comparer(proto.Equal),
		// Use big.Rat.Cmp for big.Rats
		cmp.Comparer(func(x, y *big.Rat) bool {
			if x == nil || y == nil {
				return x == y
			}
			return x.Cmp(y) == 0
		}),
		// NaNs compare equal
		cmp.FilterValues(func(x, y float64) bool {
			return math.IsNaN(x) && math.IsNaN(y)
		}, alwaysEqual),
		cmp.FilterValues(func(x, y float32) bool {
			return math.IsNaN(float64(x)) && math.IsNaN(float64(y))
		}, alwaysEqual),
	}
)

// Equal tests two values for equality.
func Equal(x, y interface{}, opts ...cmp.Option) bool {
	// Put default options at the end. Order doesn't matter.
	opts = append(opts[:len(opts):len(opts)], defaultCmpOptions...)
	return cmp.Equal(x, y, opts...)
}

// Diff reports the differences between two values.
// Diff(x, y) == "" iff Equal(x, y).
func Diff(x, y interface{}, opts ...cmp.Option) string {
	// Put default options at the end. Order doesn't matter.
	opts = append(opts[:len(opts):len(opts)], defaultCmpOptions...)
	return cmp.Diff(x, y, opts...)
}

// IsDeepEqualTest checks if the source and comparison are deeply equal and reports any errors to the testing.T instance.
func IsDeepEqualTest(t *testing.T, source interface{}, comparison interface{}, category string, description string) bool {
	equal, err := IsDeepEqual(source, comparison, category, description)
	if err != nil {
		t.Errorf("%v() test \"%v\" output does not match expected data: %v", category, description, err)
		return false
	}
	if !equal {
		t.Errorf("%v() test \"%v\" output does not match expected data", category, description)
	}
	return equal
}

// IsDeepEqual performs both IsMapVersionDeepEqual and IsStringVersionDeepEqual checks
func IsDeepEqual(source interface{}, comparison interface{}, category string, description string) (bool, error) {
	equal := IsMapVersionDeepEqual(source, comparison, category, description)
	if !equal {
		return false, nil
	}

	equal, err := IsStringVersionDeepEqual(source, comparison, category, description)
	if err != nil {
		return false, err
	}
	return equal, nil
}

// IsMapVersionDeepEqual checks if the result is equal to the expected result
func IsMapVersionDeepEqual(source interface{}, comparison interface{}, category string, description string) bool {
	// Convert to maps if not already maps
	expectedMap := asMap(comparison)
	resultMap := asMap(source)

	if !reflect.DeepEqual(resultMap, expectedMap) {
		return false
	}
	return true
}

// asMap converts a value to a map representation.
func asMap(v interface{}) interface{} {
	if v == nil {
		return nil
	}

	t := reflect.TypeOf(v)
	if t != nil && t.Kind() == reflect.Map {
		return v
	}
	return structs.Map(v)
}
