# firestruct
[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/bennovw/firestruct)
[![Codecov](https://codecov.io/gh/bennovw/firestruct/branch/main/graph/badge.svg?token=MDBGUOQY6P)](https://codecov.io/gh/bennovw/firestruct)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg?style=flat)](https://goreportcard.com/report/github.com/bennovw/firestruct)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbennovw%2Ffirestruct.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbennovw%2Ffirestruct?ref=badge_shield)

This package deserializes protojson encoded Firestore documents as found in Firestore Cloud Events into a native Go map[string]interface{} or struct (without any protojson tags).

## Why Should You Use This Package?
This package is useful when you need to read incoming Firestore documents within a Google Cloud Function written in Go. As of writing, the official Firestore Go SDK is not compatible with Cloud Functions because it does not support Firestore triggers or Cloud Events. This unofficial Go SDK for Firestore Cloud Events implements the missing functionality.

Google Cloud Functions can be triggered by events in Firestore (onCreate, onUpdate, onDelete, onWrite). The resulting http call to your cloud function includes the Firestore document in the http payload, but all the document's fields will be wrapped in protojson type descriptor tags. These tags makes parsing the document much more difficult than it needs to be. 

This package makes it trivial to unwrap the encoded Firestore data and unmarshal your Cloud Event into native Go data structures. You can simply re-use the same structs used to create Firestore documents to handle incoming cloud events. This in turn neatly simplifies your data processing and validation within your codebase.

## Installation
```go get github.com/bennovw/firestruct```

## Usage
See the [examples](https://github.com/bennovw/firestruct/tree/main/examples) folder for more examples.

```go
import (
    "github.com/bennovw/firestruct"
)

func MyCloudFunction(ctx context.Context, e event.Event) error {
    cloudEvent := firestruct.FirestoreCloudEvent{}
    err := json.Unmarshal(e.DataEncoded, &cloudEvent)
    if err != nil {
        fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
        return err
    }

    // Extract and unwrap a protojson encoded Firestore document from a Cloud Event
    // Outputs a flattened map[string]interface{} without Firestore protojson tags
    m, err := cloudEvent.ToMap()
    if err != nil {
        fmt.Printf("Error converting firestore document to map: %s", err)
    }

    // Unwrap and unmarshal a protojson encoded Firestore document into a struct
    x := MyStruct{}
    err = cloudEvent.DataTo(&x)
    if err != nil {
        fmt.Printf("Error converting firestore document to MyStruct: %s", err)
        return err
    }

    return nil
}

// Supports all Firestore data types, including nested maps and arrays,
// Firestore struct tags are optional
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
```

## Advanced Example
The package also provides two stand-alone functions to flatten a subset of Firestore data into a map[string]interface{} or unmarshal data directly into a struct without having to rely on type assertions or json.Marshal followed by json.Unmarshal
```go
import (
    "github.com/bennovw/firestruct"
)

func MyCloudFunction(ctx context.Context, e event.Event) error {
    cloudEvent := firestruct.FirestoreCloudEvent{}
    err := json.Unmarshal(e.DataEncoded, &cloudEvent)
    if err != nil {
        fmt.Printf("Error unmarshalling firestore cloud event: %s", err)
        return err
    }

    // Unwraps a protojson encoded Firestore document, outputs a flattened map[string]interface{}
    uf, err := firestruct.UnwrapFirestoreFields(cloudEvent.Value.Fields)
    if err != nil {
        fmt.Printf("Error unwrapping firestore data: %s", err)
    }

    // Unmarshals a map[string]interface{} directly into a struct
    st := MyStruct{}
    err = firestruct.DataTo(&st, uf)
    if err != nil {
        fmt.Printf("Error populating MyStruct: %s", err)
    }

    return nil
}
```

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbennovw%2Ffirestruct.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbennovw%2Ffirestruct?ref=badge_large)
