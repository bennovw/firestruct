package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/bennovw/firestruct"

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

// MyCloudFunction is an example of how to use the firestruct package in a Google Cloud Function
// The cloud function would be triggered by a Cloud Event after a Firestore Document changed in a collection
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
		fmt.Printf("Error populating MyStruct: %s", err)
		return err
	}

	// Do something with x

	return nil
}

func main() {
	// All fields in the Firestore document will be wrapped by protojson tags.
	// This is just an offline example, you would normally receive the Cloud Event within a Google Cloud Function or other Cloud Event listener after a Cloud Firestore trigger is activated.
	cloudEvent := firestruct.FirestoreCloudEvent{}
	err := json.Unmarshal([]byte(firestoreCloudEventJSON), &cloudEvent)
	if err != nil {
		fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
	}

	// Extract and unwrap a protojson encoded Firestore document from a Cloud Event
    // Outputs a flattened map[string]interface{} without Firestore protojson tags
	m, err := cloudEvent.ToMap()
	if err != nil {
		fmt.Printf("Error converting firestore document to map: %s", err)
	}
	spew.Dump(m)

	// Unwrap and unmarshal a protojson encoded Firestore document into a struct
	s := MyStruct{}
	err = cloudEvent.DataTo(&s)
	if err != nil {
		fmt.Printf("Error converting firestore document to MyStruct: %s", err)
	}
	//spew.Dump(s)

	// Unwraps a protojson encoded Firestore document, outputs a flattened map[string]interface{}
	uf, err := firestruct.UnwrapFirestoreFields(cloudEvent.Value.Fields)
	if err != nil {
		fmt.Printf("Error unwrapping firestore data: %s", err)
	}
	//spew.Dump(uf)

    // Unmarshals a map[string]interface{} directly into a struct
	st := MyStruct{}
	err = firestruct.DataTo(&st, uf)
	if err != nil {
		fmt.Printf("Error converting reflect.pointer to MyStruct: %s", err)
	}
	//spew.Dump(st)
}
