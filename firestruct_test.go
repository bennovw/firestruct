package firestruct

import (
	"encoding/json"
	"testing"

	"github.com/bennovw/firestruct/internal/testutil"
)

func TestFirestoreCloudEvent(t *testing.T) {
	var testEvent = testutil.TestFirebaseCloudEvents[0]
	var testCloudEventBytes, _ = json.Marshal(testEvent)
	var testCloudEvent = []byte(testCloudEventBytes)
	var FirestoreCloudEventTests = []testutil.TableTest{
		{
			Name:       "data to simple struct",
			TargetType: "DataToSimple",
			Input:      testCloudEvent,
			Expected:   testutil.StructResults[0],
		},
		{
			Name:       "data to tagged struct",
			TargetType: "DataToTagged",
			Input:      testCloudEvent,
			Expected:   testutil.StructResults[1],
		},
		{
			Name:       "data to flattened map",
			TargetType: "ToMap",
			Input:      testCloudEvent,
			Expected:   testutil.FlattenedMapResults[12],
		},
	}

	for _, test := range FirestoreCloudEventTests {
		thisMethodName := "FirestoreCloudEvent"

		receivedCloudEvent := FirestoreCloudEvent{}
		err := json.Unmarshal(testCloudEvent, &receivedCloudEvent)
		if err != nil {
			t.Errorf("%v() test \"%v\" returned error: %v", thisMethodName, test.Name, err)
		}

		switch test.TargetType {
		case "DataToSimple":
			var result testutil.TestSimpleStruct
			err = receivedCloudEvent.DataTo(&result)
			if err != nil {
				t.Errorf("%v() test \"%v\" returned error: %v", thisMethodName, test.Name, err)
			}
			testutil.IsDeepEqualTest(t, result, test.Expected,thisMethodName, test.Name)
		case "DataToTagged":
			var result testutil.TestTaggedStruct
			err = receivedCloudEvent.DataTo(&result)
			if err != nil {
				t.Errorf("%v() test \"%v\" returned error: %v", thisMethodName, test.Name, err)
			}
			testutil.IsDeepEqualTest(t, result, test.Expected,thisMethodName, test.Name)
		case "ToMap":
			result, err := receivedCloudEvent.ToMap()
			if err != nil {
				t.Errorf("%v() test \"%v\" returned error: %v", thisMethodName, test.Name, err)
			}
			testutil.IsDeepEqualTest(t, result, test.Expected,thisMethodName, test.Name)
		default:
			t.Errorf("%v() test \"%v\" method not covered", thisMethodName, test.Name)
		}
	}
}
