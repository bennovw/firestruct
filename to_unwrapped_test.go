package firestruct

import (
	"testing"

	"github.com/bennovw/firestruct/internal/testutil"
)

type UnwrappedTableTest struct {
	Name     string
	Input    map[string]any
	Expected map[string]any
}

var firestoreUnwrapTests = []UnwrappedTableTest{
	{
		Name:     "Firestore Time",
		Input:    testutil.TestFirebaseDocFields[0],
		Expected: testutil.FlattenedMapResults[0],
	},
	{
		Name:     "Firestore String",
		Input:    testutil.TestFirebaseDocFields[1],
		Expected: testutil.FlattenedMapResults[1],
	},
	{
		Name:     "Firestore UUID",
		Input:    testutil.TestFirebaseDocFields[2],
		Expected: testutil.FlattenedMapResults[2],
	},
	{
		Name:     "Firestore bool",
		Input:    testutil.TestFirebaseDocFields[3],
		Expected: testutil.FlattenedMapResults[3],
	},
	{
		Name:     "Firestore int",
		Input:    testutil.TestFirebaseDocFields[4],
		Expected: testutil.FlattenedMapResults[4],
	},
	{
		Name:     "Firestore double",
		Input:    testutil.TestFirebaseDocFields[5],
		Expected: testutil.FlattenedMapResults[5],
	},
	{
		Name:     "Firestore bytes",
		Input:    testutil.TestFirebaseDocFields[6],
		Expected: testutil.FlattenedMapResults[6],
	},
	{
		Name:     "Firestore nil",
		Input:    testutil.TestFirebaseDocFields[7],
		Expected: testutil.FlattenedMapResults[7],
	},
	{
		Name:     "Firestore reference",
		Input:    testutil.TestFirebaseDocFields[8],
		Expected: testutil.FlattenedMapResults[8],
	},
	{
		Name:     "Firestore geopoint",
		Input:    testutil.TestFirebaseDocFields[9],
		Expected: testutil.FlattenedMapResults[9],
	},
	{
		Name:     "Firestore Map",
		Input:    testutil.TestFirebaseDocFields[10],
		Expected: testutil.FlattenedMapResults[10],
	},
	{
		Name:     "Firestore Array",
		Input:    testutil.TestFirebaseDocFields[11],
		Expected: testutil.FlattenedMapResults[11],
	},
	{
		Name:     "Firestore Nested Fields",
		Input:    testutil.TestFirebaseDocFields[12],
		Expected: testutil.FlattenedMapResults[12],
	},
}

func TestUnwrapFirestoreFields(t *testing.T) {
	thisFunctionName := "UnwrapFirestoreFields"
	for _, test := range firestoreUnwrapTests {
		t.Run(test.Name, func(t *testing.T) {
			result, err := UnwrapFirestoreFields(test.Input)
			if err != nil {
				t.Errorf("%v() test \"%v\" returned error: %v", thisFunctionName, test.Name, err)
			}

			testutil.IsDeepEqual(t, thisFunctionName, test.Name, result, test.Expected)
		})
	}

}
