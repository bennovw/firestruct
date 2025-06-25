package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudevents/sdk-go/v2/event"

	"github.com/bennovw/firestruct"
	"github.com/bennovw/firestruct/internal/testutil"
)

// IntegrationTest tests the firestruct package's ability to decode Firestore cloud events containing protojson data.
// This Google Cloud Function can be triggered by a firestore document modification event,
// it will then compare the decoded event data against expected results defined in testutil.StructResults
func IntegrationTest(ctx context.Context, e event.Event) error {
	log.Printf("Received Cloud Event: %s", e.Type())
	log.Printf("Event ID: %s", e.ID())
	log.Printf("Event Source: %s", e.Source())

	// Parse the Cloud Event data into a FirestoreCloudEvent
	cloudEvent := firestruct.FirestoreCloudEvent{}
	err := json.Unmarshal(e.DataEncoded, &cloudEvent)
	if err != nil {
		return fmt.Errorf("error unmarshalling firestore cloud event: %w", err)
	}

	// Extract the document data and decode it as TestTaggedStruct
	decoded := testutil.TestTaggedStruct{}
	err = cloudEvent.DataTo(&decoded)
	if err != nil {
		return fmt.Errorf("error converting firestore document to TestTaggedStruct: %w", err)
	}

	// Log the decoded data for verification
	log.Printf("Successfully decoded Firebase document to TestTaggedStruct:")
	log.Printf("  Time: %v", decoded.Time)
	log.Printf("  String: %s", decoded.String)
	log.Printf("  UUID: %s", decoded.UUID)
	log.Printf("  Bool: %t", decoded.Bool)
	log.Printf("  Int: %d", decoded.Int)
	log.Printf("  Double: %f", decoded.Double)
	log.Printf("  Bytes: %s", string(decoded.Bytes))
	log.Printf("  Nil: %v", decoded.Nil)
	log.Printf("  GeoPoint: lat=%f, lng=%f", decoded.GeoPoint.Latitude, decoded.GeoPoint.Longitude)
	log.Printf("  Ref: %s", decoded.Ref)
	log.Printf("  NestedMap keys: %d", len(decoded.NestedMap))

	// Compare the result with the expected data in StructResults[1] using testutil.IsDeepEqual()
	expectedData := testutil.StructResults[1]
	equal, err := testutil.IsDeepEqual(decoded, expectedData, "FirebaseDocumentModifiedFunction", "TestTaggedStruct comparison")
	if !equal {
		return fmt.Errorf("firebase event payload does not match test data: %v", err)
	} else {
		log.Printf("firebase event payload matches expected test data")
	}

	return nil
}
