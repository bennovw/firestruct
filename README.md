# firestruct
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/bennovw/firestruct)
[![codecov](https://codecov.io/gh/bennovw/firestruct/branch/master/graph/badge.svg)](https://codecov.io/gh/bennovw/firestruct)
[![Go Report Card](https://goreportcard.com/badge/github.com/bennovw/firestruct)](https://goreportcard.com/report/github.com/bennovw/firestruct)

This package removes Firestore protojson tags and populates Go structs.

The package can be used by cloud functions to handle Firestore Cloud Events containing encoded Firestore documents.

## Why Should I Use This Package?
Firestore Cloud Events contain JSON encoded Firestore documents. The encoded Firestore document fields are all wrapped by protojson type descriptor tags. These tags are not needed when handling Cloud Events, and make it harder to work with Firestore data in Go cloud functions. This package can remove the Firestore protojson tags and populate Go structs using native Go types. As a result, you can easily re-use your existing models in your cloud functions.

## Simple Usage
See the [examples] folder (https://github.com/bennovw/firestruct/examples) for all examples.

```go
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

func MyCloudFunction(ctx context.Context, e event.Event) error {
	cloudEvent := firestruct.FirestoreCloudEvent{}
	err := json.Unmarshal(e.DataEncoded, &cloudEvent)
	if err != nil {
		fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
		return err
	}

    // DataTo() converts the Firestore document into a struct, using the struct tags to map the Firestore document fields to the struct fields.
	// This method is available to both FirestoreDocument and FirestoreCloudEvent types.
	x := MyStruct{}
	err = cloudEvent.DataTo(&x)
	if err != nil {
		fmt.Printf("Error converting firestore document to MyStruct: %s", err)
		return err
	}

	// Do something with x

    // ToMap() removes any Firestore protojson tags by converting the Cloud Event to a map[string]interface{}
	// This method is available to both FirestoreDocument and FirestoreCloudEvent types.
    m, err := cloudEvent.ToMap()
	if err != nil {
		fmt.Printf("Error converting firestore document to map: %s", err)
	}

    // Do something with m

	return nil
}
```

## Advanced Usage
For advanced use cases, the package provides helper functions to unwrap Firestore protojson tags from any map[string]interface{} or populate a struct using any type of source data.
```go
func MyCloudFunction(ctx context.Context, e event.Event) error {
	cloudEvent := firestruct.FirestoreCloudEvent{}
	err := json.Unmarshal(e.DataEncoded, &cloudEvent)
	if err != nil {
		fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
		return err
	}

	// For advanced use cases, UnwrapFirestoreFields() will remove any Firestore protojson tags form a map[string]interface{}
	uf, err := firestruct.UnwrapFirestoreFields(cloudEvent.Value.Fields)
	if err != nil {
		fmt.Printf("Error unwrapping firestore data: %s", err)
	}
	
    // Do something with uf

	// For advanced use cases, DataTo() will populate a struct using any type of source data.
	st := MyStruct{}
	err = firestruct.DataTo(&st, uf)
	if err != nil {
		fmt.Printf("Error converting reflect.pointer to MyStruct: %s", err)
	}
	
    // Do something with st
    return nil
}
```

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbennovw%2Ffirestruct.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbennovw%2Ffirestruct?ref=badge_large)
