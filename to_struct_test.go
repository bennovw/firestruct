package firestruct

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/latlng"

	"github.com/bennovw/firestruct/internal/testutil"
)

var firestoreToStructTests = []testutil.TableTest{
	{
		Name:       "data to simple struct",
		TargetType: "SimpleStruct",
		Input:      testutil.TestFirebaseDocFields[12],
		Expected:   testutil.StructResults[0],
	},
	{
		Name:       "data to tagged struct",
		TargetType: "TaggedStruct",
		Input:      testutil.TestFirebaseDocFields[12],
		Expected:   testutil.StructResults[1],
	},
}

func BenchmarkTestDataToReflectPointerSimple(b *testing.B) {
	benchmarkName := "BenchmarkTestDataToReflectPointerSimple"
	testStruct := firestoreToStructTests[0]
	testData, err := UnwrapFirestoreFields(testStruct.Input.(map[string]interface{}))
	if err != nil {
		b.Errorf("%v() test \"%v\" returned error running UnwrapFirestoreFields(): %v", benchmarkName, testStruct.Name, err)
	}

	for i := 0; i < b.N; i++ {
		b.Run(benchmarkName, func(b *testing.B) {
			result := testutil.TestSimpleStruct{}
			_ = dataToReflectPointer(reflect.ValueOf(&result).Elem(), testData)
		})
	}
}

func BenchmarkTestDataToReflectPointerTagged(b *testing.B) {
	benchmarkName := "BenchmarkTestDataToReflectPointerTagged"
	testStruct := firestoreToStructTests[1]
	testData, err := UnwrapFirestoreFields(testStruct.Input.(map[string]interface{}))
	if err != nil {
		b.Errorf("%v() test \"%v\" returned error running UnwrapFirestoreFields(): %v", benchmarkName, testStruct.Name, err)
	}

	for i := 0; i < b.N; i++ {
		b.Run(benchmarkName, func(b *testing.B) {
			result := testutil.TestTaggedStruct{}
			_ = dataToReflectPointer(reflect.ValueOf(&result).Elem(), testData)
		})
	}
}

func TestDataToReflectPointer(t *testing.T) {
	thisFunctionName := "dataToReflectPointer"

	for _, test := range firestoreToStructTests {
		t.Run(test.Name, func(t *testing.T) {
			testInput, ok := test.Input.(map[string]interface{})
			if !ok {
				t.Errorf("%v() test \"%v\" test data invalid, Input1 is not a map[string]interface{}", thisFunctionName, test.Name)
			}

			unwrapped, err := UnwrapFirestoreFields(testInput)
			if err != nil {
				t.Errorf("%v() test \"%v\" returned error running UnwrapFirestoreFields(): %v", thisFunctionName, test.Name, err)
			}

			switch test.Expected.(type) {
			case testutil.TestSimpleStruct:
				result := testutil.TestSimpleStruct{}
				err = dataToReflectPointer(reflect.ValueOf(&result).Elem(), unwrapped)
				if err != nil {
					t.Errorf("%v() test \"%v\" returned error: %v", thisFunctionName, test.Name, err)
				}

				testutil.IsDeepEqualTest(t, result, test.Expected, thisFunctionName, test.Name)
			case testutil.TestTaggedStruct:
				result := testutil.TestTaggedStruct{}
				err = dataToReflectPointer(reflect.ValueOf(&result).Elem(), unwrapped)
				if err != nil {
					t.Errorf("%v() test \"%v\" returned error: %v", thisFunctionName, test.Name, err)
				}

				testutil.IsDeepEqualTest(t, result, test.Expected, thisFunctionName, test.Name)
			default:
				t.Errorf("%v() test \"%v\" expected result data type is not covered", thisFunctionName, test.Name)
			}
		})
	}

}

func TestDataToReflectPointerTypes(t *testing.T) {
	type inner struct {
		N int
	}
	var heldPtr any = new(int)
	testTime := time.Date(2023, 5, 1, 12, 30, 0, 0, time.UTC)
	testUUID := uuid.MustParse("b4c5d6e7-1234-4abc-8def-0123456789ab")

	tests := []struct {
		name    string
		target  any // pointer to the value being populated
		data    any
		want    any // expected value of *target
		wantErr bool
	}{
		// nil
		{name: "nil into pointer", target: func() any { p := new(int); return &p }(), data: nil, want: (*int)(nil)},
		{name: "nil into int has no effect", target: func() any { i := 7; return &i }(), data: nil, want: 7},

		// []byte
		{name: "bytes from base64 string", target: new([]byte), data: "aGVsbG8=", want: []byte("hello")},
		{name: "bytes from bytes", target: new([]byte), data: []byte("hi"), want: []byte("hi")},
		{name: "bytes from invalid base64", target: new([]byte), data: "!!!", wantErr: true},
		{name: "bytes from int", target: new([]byte), data: 1, wantErr: true},

		// time.Time
		{name: "time from time", target: new(time.Time), data: testTime, want: testTime},
		{name: "time from RFC3339 string", target: new(time.Time), data: "2023-05-01T12:30:00Z", want: testTime},
		{name: "time from invalid string", target: new(time.Time), data: "yesterday", wantErr: true},
		{name: "time from int", target: new(time.Time), data: 1, wantErr: true},

		// latlng.LatLng
		{name: "latlng from map", target: new(latlng.LatLng), data: map[string]any{"latitude": 1.5, "longitude": -2.5}, want: latlng.LatLng{Latitude: 1.5, Longitude: -2.5}},
		{name: "latlng with invalid latitude", target: new(latlng.LatLng), data: map[string]any{"latitude": "1", "longitude": 2.0}, wantErr: true},
		{name: "latlng with invalid longitude", target: new(latlng.LatLng), data: map[string]any{"latitude": 1.0, "longitude": "2"}, wantErr: true},
		{name: "latlng from string", target: new(latlng.LatLng), data: "1,2", wantErr: true},

		// uuid.UUID
		{name: "uuid from string", target: new(uuid.UUID), data: testUUID.String(), want: testUUID},
		{name: "uuid from invalid string", target: new(uuid.UUID), data: "not-a-uuid", wantErr: true},
		{name: "uuid from int", target: new(uuid.UUID), data: 1, wantErr: true},

		// bool and string
		{name: "bool", target: new(bool), data: true, want: true},
		{name: "bool from string", target: new(bool), data: "true", wantErr: true},
		{name: "string", target: new(string), data: "hi", want: "hi"},
		{name: "string from int", target: new(string), data: 1, wantErr: true},

		// signed integers
		{name: "int from int", target: new(int), data: 42, want: 42},
		{name: "int from int8", target: new(int), data: int8(8), want: 8},
		{name: "int from int16", target: new(int), data: int16(16), want: 16},
		{name: "int from int32", target: new(int), data: int32(32), want: 32},
		{name: "int64 from int64", target: new(int64), data: int64(64), want: int64(64)},
		{name: "int from float64", target: new(int), data: 3.0, want: 3},
		{name: "int8 overflow", target: new(int8), data: 300, wantErr: true},
		{name: "int from string", target: new(int), data: "42", wantErr: true},

		// unsigned integers
		{name: "uint8 from uint8", target: new(uint8), data: uint8(8), want: uint8(8)},
		{name: "uint16 from uint16", target: new(uint16), data: uint16(16), want: uint16(16)},
		{name: "uint32 from uint32", target: new(uint32), data: uint32(32), want: uint32(32)},
		{name: "uint32 from uint64", target: new(uint32), data: uint64(64), want: uint32(64)},
		{name: "uint8 overflow", target: new(uint8), data: uint16(300), wantErr: true},
		{name: "uint8 from int", target: new(uint8), data: 1, wantErr: true},

		// floats
		{name: "float64 from float64", target: new(float64), data: 1.25, want: 1.25},
		{name: "float64 from float32", target: new(float64), data: float32(1.5), want: 1.5},
		{name: "float32 overflow", target: new(float32), data: 1e300, wantErr: true},
		{name: "float64 from int", target: new(float64), data: 1, want: 1.0},

		// slices
		{name: "slice grows", target: new([]string), data: []any{"a", "b"}, want: []string{"a", "b"}},
		{name: "slice shrinks", target: &[]string{"x", "y", "z"}, data: []any{"a"}, want: []string{"a"}},
		{name: "slice from map", target: new([]string), data: map[string]any{}, wantErr: true},
		{name: "slice with invalid element", target: new([]string), data: []any{"a", 1}, wantErr: true},

		// arrays
		{name: "array zeroes extra elements", target: &[3]int{9, 9, 9}, data: []any{1}, want: [3]int{1, 0, 0}},
		{name: "array ignores extra values", target: new([2]int), data: []any{1, 2, 3}, want: [2]int{1, 2}},
		{name: "array from string", target: new([2]int), data: "12", wantErr: true},

		// maps
		{name: "map", target: new(map[string]int), data: map[string]any{"a": 1}, want: map[string]int{"a": 1}},
		{name: "map overwrites existing keys", target: &map[string]int{"a": 0, "b": 2}, data: map[string]any{"a": 1}, want: map[string]int{"a": 1, "b": 2}},
		{name: "map with non-string keys", target: new(map[int]int), data: map[string]any{"1": 1}, wantErr: true},
		{name: "map with invalid value", target: new(map[string]int), data: map[string]any{"a": "1"}, wantErr: true},
		{name: "map from slice", target: new(map[string]int), data: []any{1}, wantErr: true},

		// pointers
		{name: "nil pointer is allocated", target: new(*int), data: 5, want: func() *int { i := 5; return &i }()},

		// structs
		{name: "struct", target: new(inner), data: map[string]any{"N": 1, "Unknown": 2}, want: inner{N: 1}},
		{name: "struct with invalid field", target: new(inner), data: map[string]any{"N": "1"}, wantErr: true},
		{name: "struct from slice", target: new(inner), data: []any{1}, wantErr: true},

		// interfaces
		{name: "empty interface", target: new(any), data: "hi", want: "hi"},
		{name: "empty interface holding a pointer", target: &heldPtr, data: 5, want: any(func() *int { i := 5; return &i }())},
		{name: "non-empty interface", target: new(fmt.Stringer), data: "hi", wantErr: true},

		// unsupported kinds
		{name: "channel", target: new(chan int), data: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.target).Elem()
			err := dataToReflectPointer(v, tt.data)
			if tt.wantErr {
				if err == nil {
					t.Errorf("dataToReflectPointer() = nil error, want error (value: %v)", v.Interface())
				}
				return
			}
			if err != nil {
				t.Fatalf("dataToReflectPointer() returned error: %v", err)
			}
			if got := v.Interface(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("dataToReflectPointer() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestPopulateStructErrorIncludesField(t *testing.T) {
	type inner struct {
		Count int
	}
	err := dataToReflectPointer(reflect.ValueOf(new(inner)).Elem(), map[string]any{"Count": "x"})
	if err == nil || !strings.Contains(err.Error(), "Count") {
		t.Errorf("dataToReflectPointer() error = %v, want error mentioning field Count", err)
	}
}

func TestOverflowErr(t *testing.T) {
	err := dataToReflectPointer(reflect.ValueOf(new(int8)).Elem(), 300)
	want := "value 300 overflows type int8"
	if err == nil || err.Error() != want {
		t.Errorf("dataToReflectPointer() error = %v, want %q", err, want)
	}
}
