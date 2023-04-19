package firestruct

import (
	"reflect"
	"testing"

	"github.com/bennovw/firestruct/internal/testutil"

	"github.com/davecgh/go-spew/spew"
)

type UnwrappedTableTest struct {
	Name     string
	Input    map[string]any
	Expected map[string]any
}

var firestoreUnwrapTests = []UnwrappedTableTest{
	{
		Name:     "Firestore Time",
		Input:    testutil.TestFirebaseDocs[0],
		Expected: testutil.ResultFirebaseDocs[0],
	},
	{
		Name:     "Firestore String",
		Input:    testutil.TestFirebaseDocs[1],
		Expected: testutil.ResultFirebaseDocs[1],
	},
	{
		Name:     "Firestore UUID",
		Input:    testutil.TestFirebaseDocs[2],
		Expected: testutil.ResultFirebaseDocs[2],
	},
	{
		Name:     "Firestore bool",
		Input:    testutil.TestFirebaseDocs[3],
		Expected: testutil.ResultFirebaseDocs[3],
	},
	{
		Name:     "Firestore int",
		Input:    testutil.TestFirebaseDocs[4],
		Expected: testutil.ResultFirebaseDocs[4],
	},
	{
		Name:     "Firestore double",
		Input:    testutil.TestFirebaseDocs[5],
		Expected: testutil.ResultFirebaseDocs[5],
	},
	{
		Name:     "Firestore bytes",
		Input:    testutil.TestFirebaseDocs[6],
		Expected: testutil.ResultFirebaseDocs[6],
	},
	{
		Name:     "Firestore nil",
		Input:    testutil.TestFirebaseDocs[7],
		Expected: testutil.ResultFirebaseDocs[7],
	},
	{
		Name:     "Firestore reference",
		Input:    testutil.TestFirebaseDocs[8],
		Expected: testutil.ResultFirebaseDocs[8],
	},
	{
		Name:     "Firestore geopoint",
		Input:    testutil.TestFirebaseDocs[9],
		Expected: testutil.ResultFirebaseDocs[9],
	},
	{
		Name:     "Firestore Map",
		Input:    testutil.TestFirebaseDocs[10],
		Expected: testutil.ResultFirebaseDocs[10],
	},
	{
		Name:     "Firestore Array",
		Input:    testutil.TestFirebaseDocs[11],
		Expected: testutil.ResultFirebaseDocs[11],
	},
	{
		Name:     "Firestore Nested Fields",
		Input:    testutil.TestFirebaseDocs[12],
		Expected: testutil.ResultFirebaseDocs[12],
	},
}

func TestUnwrapFirestoreFields(t *testing.T) {
	for _, test := range firestoreUnwrapTests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := UnwrapFirestoreFields(test.Input)
			if !reflect.DeepEqual(result, test.Expected) {
				spew.Dump(result, test.Expected)
				t.Errorf("UnwrapFirestoreFields() test \"%v\" Result: %v Want: %v", test.Name, result, test.Expected)
			}
			if err != nil {
				t.Errorf("UnwrapFirestoreFields() test \"%v\" returned error: %v", test.Name, err)
			}
		})
	}

}
