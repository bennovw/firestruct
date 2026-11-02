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

func TestFirestoreCloudEventDocument(t *testing.T) {
	e := FirestoreCloudEvent{Value: FirestoreDocument{Name: "doc"}}
	if got := e.Document(); got != &e.Value {
		t.Errorf("Document() = %p, want pointer to Value (%p)", got, &e.Value)
	}
}

func TestDataToErrors(t *testing.T) {
	var nilPtr *struct{}
	tests := []struct {
		name   string
		target any
	}{
		{name: "nil", target: nil},
		{name: "non-pointer", target: struct{}{}},
		{name: "nil pointer", target: nilPtr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DataTo(tt.target, map[string]any{}); err == nil {
				t.Errorf("DataTo(%v) = nil, want error", tt.target)
			}
		})
	}
}

func TestDataToMap(t *testing.T) {
	data := map[string]any{"a": 1}
	got := map[string]any{}
	if err := DataTo(&got, data); err != nil {
		t.Fatalf("DataTo() returned error: %v", err)
	}
	if got["a"] != 1 {
		t.Errorf("DataTo() = %v, want %v", got, data)
	}
}

func TestFirestoreDocumentErrors(t *testing.T) {
	var nilDoc *FirestoreDocument
	if _, err := nilDoc.ToMap(); err == nil {
		t.Error("ToMap() on nil document = nil error, want error")
	}

	invalid := FirestoreDocument{Fields: map[string]any{"a": "not wrapped"}}
	if _, err := invalid.ToMap(); err == nil {
		t.Error("ToMap() with invalid fields = nil error, want error")
	}
	var target struct{ A string }
	if err := invalid.DataTo(&target); err == nil {
		t.Error("DataTo() with invalid fields = nil error, want error")
	}

	e := FirestoreCloudEvent{Value: invalid}
	if _, err := e.ToMap(); err == nil {
		t.Error("FirestoreCloudEvent.ToMap() with invalid fields = nil error, want error")
	}
	if err := e.DataTo(&target); err == nil {
		t.Error("FirestoreCloudEvent.DataTo() with invalid fields = nil error, want error")
	}
}
