package firestruct

import (
	"reflect"
	"testing"

	"github.com/bennovw/firestruct/internal/testutil"

	"github.com/davecgh/go-spew/spew"
)

var firestoreToStructTests = []testutil.TableTest{
	{
		Name:     "data to simple struct",
		Input:    testutil.TestFirebaseDocs[12],
		Expected: testutil.ResultStructs[0],
	},
	{
		Name:     "data to tagged struct",
		Input:    testutil.TestFirebaseDocs[12],
		Expected: testutil.ResultStructs[1],
	},
}

func TestDataToReflectPointer(t *testing.T) {
	testDocSimple := testutil.TestDocSimple{}
	testDocTagged := testutil.TestDocTagged{}

	for _, test := range firestoreToStructTests {
		t.Run(test.Name, func(t *testing.T) {
			x, ok := test.Input.(map[string]interface{})
			if !ok {
				t.Errorf("dataToReflectPointer() test \"%v\" test data invalid, Input1 is not a map[string]interface{}", test.Name)
			}

			unwrapped, err := UnwrapFirestoreFields(x)
			if err != nil {
				t.Errorf("dataToReflectPointer() test \"%v\" returned error running UnwrapFirestoreFields(): %v", test.Name, err)
			}

			switch test.Expected.(type) {
			case testutil.TestDocSimple:
				err = dataToReflectPointer(reflect.ValueOf(&testDocSimple).Elem(), unwrapped)
				if err != nil {
					t.Errorf("dataToReflectPointer() test \"%v\" returned error: %v", test.Name, err)
				}

				CheckResultEquality(t, test.Name, testDocSimple, test.Expected)
			case testutil.TestDocTagged:
				err = dataToReflectPointer(reflect.ValueOf(&testDocTagged).Elem(), unwrapped)
				if err != nil {
					t.Errorf("dataToReflectPointer() test \"%v\" returned error: %v", test.Name, err)
				}

				CheckResultEquality(t, test.Name, testDocTagged, test.Expected)
			default:
				t.Errorf("dataToReflectPointer() test \"%v\" expected result data type is not covered", test.Name)
			}
		})
	}

}

// CheckResultEquality checks if the result is equal to the expected result
func CheckResultEquality(t *testing.T, name string, result interface{}, expected interface{}) {
	if !reflect.DeepEqual(result, expected) {
		spew.Dump(result, expected)
		t.Errorf("dataToReflectPointer() test \"%v\" Result: %v Want: %v", name, result, expected)
	}
}
