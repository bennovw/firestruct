package firestruct

import (
	"reflect"
	"testing"
	"time"

	"google.golang.org/genproto/googleapis/type/latlng"

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

			testutil.IsDeepEqualTest(t, result, test.Expected, thisFunctionName, test.Name)
		})
	}

}

func TestUnwrapDouble(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    float64
		wantErr bool
	}{
		{name: "float64", input: 1.5, want: 1.5},
		{name: "float32", input: float32(2.5), want: 2.5},
		{name: "int", input: 3, want: 3},
		{name: "int32", input: int32(4), want: 4},
		{name: "int64", input: int64(5), want: 5},
		{name: "string", input: "6.75", want: 6.75},
		{name: "invalid string", input: "six", wantErr: true},
		{name: "unsupported type", input: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unwrapDouble(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unwrapDouble(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("unwrapDouble(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnwrapInt(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    int
		wantErr bool
	}{
		{name: "string", input: "42", want: 42},
		{name: "invalid string", input: "forty-two", wantErr: true},
		{name: "float32", input: float32(1), want: 1},
		{name: "float64", input: 2.0, want: 2},
		{name: "int", input: 3, want: 3},
		{name: "int32", input: int32(4), want: 4},
		{name: "int64", input: int64(5), want: 5},
		{name: "unsupported type", input: true, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unwrapInt(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unwrapInt(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("unwrapInt(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnwrapBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{name: "base64 string", input: "aGVsbG8=", want: []byte("hello")},
		{name: "bytes", input: []byte("hi"), want: []byte("hi")},
		{name: "nil", input: nil, wantErr: true},
		{name: "invalid base64", input: "!!!", wantErr: true},
		{name: "unsupported type", input: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unwrapBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unwrapBytes(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unwrapBytes(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnwrapGeoPoint(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    latlng.LatLng
		wantErr bool
	}{
		{name: "nil", input: nil},
		{name: "map", input: map[string]any{"latitude": 1.5, "longitude": 2.5}, want: latlng.LatLng{Latitude: 1.5, Longitude: 2.5}},
		{name: "invalid latitude", input: map[string]any{"latitude": "1", "longitude": 2.5}, wantErr: true},
		{name: "invalid longitude", input: map[string]any{"latitude": 1.5, "longitude": "2"}, wantErr: true},
		{name: "unsupported type", input: "1,2", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unwrapGeoPoint(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unwrapGeoPoint(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got.Latitude != tt.want.Latitude || got.Longitude != tt.want.Longitude {
				t.Errorf("unwrapGeoPoint(%v) = %v, want %v", tt.input, got.String(), tt.want.String())
			}
		})
	}
}

func TestUnwrapTimestamp(t *testing.T) {
	testTime := time.Date(2023, 5, 1, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name    string
		input   any
		want    time.Time
		wantErr bool
	}{
		{name: "nil", input: nil},
		{name: "time", input: testTime, want: testTime},
		{name: "RFC3339 string", input: "2023-05-01T12:30:00Z", want: testTime},
		{name: "invalid string", input: "yesterday", wantErr: true},
		{name: "unsupported type", input: 1, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := unwrapTimestamp(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("unwrapTimestamp(%v) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if !got.Equal(tt.want) {
				t.Errorf("unwrapTimestamp(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestUnwrapErrors(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]any
	}{
		{name: "nil input", input: nil},
		{name: "non-map value", input: map[string]any{"a": "b"}},
		{name: "unsupported flat value", input: map[string]any{"a": map[string]any{"unknownValue": 1}}},
		{name: "invalid int", input: map[string]any{"a": map[string]any{protoIntTag: "x"}}},
		{name: "map without fields", input: map[string]any{"a": map[string]any{protoMapTag: map[string]any{"other": 1}}}},
		{name: "map fields not a map", input: map[string]any{"a": map[string]any{protoMapTag: map[string]any{"fields": 1}}}},
		{name: "map value not a map", input: map[string]any{"a": map[string]any{protoMapTag: 1}}},
		{name: "map with invalid nested value", input: map[string]any{"a": map[string]any{protoMapTag: map[string]any{"fields": map[string]any{"b": "c"}}}}},
		{name: "top-level map with invalid nested value", input: map[string]any{protoMapTag: map[string]any{"fields": map[string]any{"b": "c"}}}},
		{name: "array without values", input: map[string]any{"a": map[string]any{protoArrayTag: map[string]any{"other": 1}}}},
		{name: "array values not a slice", input: map[string]any{"a": map[string]any{protoArrayTag: map[string]any{"values": 1}}}},
		{name: "array value not a map", input: map[string]any{"a": map[string]any{protoArrayTag: 1}}},
		{name: "array element not a map", input: map[string]any{"a": map[string]any{protoArrayTag: map[string]any{"values": []any{1}}}}},
		{name: "array element with invalid flat value", input: map[string]any{"a": map[string]any{protoArrayTag: map[string]any{"values": []any{map[string]any{protoIntTag: "x"}}}}}},
		{name: "array element with invalid nested value", input: map[string]any{"a": map[string]any{protoArrayTag: map[string]any{"values": []any{map[string]any{"b": "c"}}}}}},
		{name: "top-level array with invalid element", input: map[string]any{protoArrayTag: map[string]any{"values": []any{1}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, err := UnwrapFirestoreFields(tt.input); err == nil {
				t.Errorf("UnwrapFirestoreFields(%v) = %v, want error", tt.input, got)
			}
		})
	}
}

func TestUnwrapEmptyContainers(t *testing.T) {
	input := map[string]any{
		"m": map[string]any{protoMapTag: map[string]any{}},
		"a": map[string]any{protoArrayTag: map[string]any{}},
	}
	got, err := UnwrapFirestoreFields(input)
	if err != nil {
		t.Fatalf("UnwrapFirestoreFields() returned error: %v", err)
	}
	if got["m"].(map[string]any) != nil || got["a"].([]any) != nil {
		t.Errorf("UnwrapFirestoreFields() = %v, want nil map and nil slice", got)
	}
}
