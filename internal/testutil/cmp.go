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

// IsDeepEqual performs both IsMapVersionDeepEqual and IsStringVersionDeepEqual checks
func IsDeepEqual(t *testing.T, funcName string, testSectionName string, result interface{}, expected interface{}) {
	IsMapVersionDeepEqual(t, funcName, testSectionName, result, expected)
	IsStringVersionDeepEqual(t, funcName, testSectionName, result, expected)
}

// IsMapVersionDeepEqual checks if the result is equal to the expected result
func IsMapVersionDeepEqual(t *testing.T, funcName string, testSectionName string, result interface{}, expected interface{}) {
	// Convert to maps if not already maps
	expectedMap := asMap(expected)
	resultMap := asMap(result)

	if !reflect.DeepEqual(resultMap, expectedMap) {
		t.Errorf("%v() test \"%v\" Result: %v Want: %v", funcName, testSectionName, result, expected)
		return
	}
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
