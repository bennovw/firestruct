package firestruct

import (
	"reflect"
	"testing"

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

				testutil.IsDeepEqual(t, thisFunctionName, test.Name, result, test.Expected)
			case testutil.TestTaggedStruct:
				result := testutil.TestTaggedStruct{}
				err = dataToReflectPointer(reflect.ValueOf(&result).Elem(), unwrapped)
				if err != nil {
					t.Errorf("%v() test \"%v\" returned error: %v", thisFunctionName, test.Name, err)
				}

				testutil.IsDeepEqual(t, thisFunctionName, test.Name, result, test.Expected)
			default:
				t.Errorf("%v() test \"%v\" expected result data type is not covered", thisFunctionName, test.Name)
			}
		})
	}

}
