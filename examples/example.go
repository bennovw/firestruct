package main

import (
	"context"
	"encoding/json"
	"firestruct"
	"fmt"
	"time"

	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/type/latlng"
)

// MyStruct is a struct that will be used to unmarshal the Firestore Document embedded in the Cloud Event
// The struct tags are used to map the Firestore Document fields to the struct fields if they are different.
type MyStruct struct {
	SomeTime   time.Time      `firestore:"timeData"`
	Title      string         `firestore:"stringData"`
	ID         uuid.UUID      `firestore:"uuidData"`
	IsWild     bool           `firestore:"boolData"`
	Age        int64          `firestore:"intData"`
	Weight     float64        `firestore:"doubleData"`
	Bytes      []byte         `firestore:"bytesData"`
	WildNull   any            `firestore:"nilData"`
	Place      latlng.LatLng  `firestore:"geoPointData"`
	NestedData map[string]any `firestore:"nestedMapData"`
}

func main() {
	// Unmarshal the Firestore Cloud Event from a JSON string.
	// All fields in the Firestore document will be wrapped by protojson tags.
	// This is just an example, you would normally receive the Cloud Event within a Google Cloud Function after a Cloud Firestore trigger is activated.
	cloudEvent := firestruct.FirestoreCloudEvent{}
	err := json.Unmarshal([]byte(firestoreCloudEventJSON), &cloudEvent)
	if err != nil {
		fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
	}

	// ToMap() removes any Firestore protojson tags by converting the Cloud Event to a map[string]interface{}
	// This method is available to both FirestoreDocument and FirestoreCloudEvent types.
	m, err := cloudEvent.ToMap()
	if err != nil {
		fmt.Printf("Error converting firestore document to map: %s", err)
	}
	spew.Dump(m)

	// DataTo() converts the Firestore document into a struct, using the struct tags to map the Firestore document fields to the struct fields.
	// This method is available to both FirestoreDocument and FirestoreCloudEvent types.
	s := MyStruct{}
	err = cloudEvent.DataTo(&s)
	if err != nil {
		fmt.Printf("Error converting firestore document to MyStruct: %s", err)
	}
	//spew.Dump(s)

	// For advanced use cases, UnwrapFirestoreFields() will remove any Firestore protojson tags form a map[string]interface{}
	uf, err := firestruct.UnwrapFirestoreFields(cloudEvent.Value.Fields)
	if err != nil {
		fmt.Printf("Error unwrapping firestore data: %s", err)
	}
	//spew.Dump(uf)

	// For advanced use cases, DataTo() will populate a struct using any type of source data.
	st := MyStruct{}
	err = firestruct.DataTo(&st, uf)
	if err != nil {
		fmt.Printf("Error converting reflect.pointer to MyStruct: %s", err)
	}
	//spew.Dump(st)
}

// MyCloudFunction is an example of how to use the firestruct package in a Google Cloud Function
// The cloud function would be triggered by a Firestore Document change and receive a Cloud Event
func MyCloudFunction(ctx context.Context, e event.Event) error {
	cloudEvent := firestruct.FirestoreCloudEvent{}
	err := json.Unmarshal(e.DataEncoded, &cloudEvent)
	if err != nil {
		fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
		return err
	}

	x := MyStruct{}
	err = cloudEvent.DataTo(&x)
	if err != nil {
		fmt.Printf("Error converting firestore document to MyStruct: %s", err)
		return err
	}

	// Do something with x

	return nil
}
