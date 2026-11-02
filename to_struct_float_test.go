package firestruct

import (
	"reflect"
	"testing"
)

func TestDataToReflectPointerFloatFromInteger(t *testing.T) {
	tests := []struct {
		name   string
		target any
		data   any
		want   any
	}{
		{name: "float64 from int", target: new(float64), data: 3, want: 3.0},
		{name: "float64 from negative int64", target: new(float64), data: int64(-64), want: -64.0},
		{name: "float32 from int32", target: new(float32), data: int32(32), want: float32(32)},
		{name: "float64 from uint", target: new(float64), data: uint(7), want: 7.0},
		{name: "float32 from uint64", target: new(float32), data: uint64(64), want: float32(64)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.target).Elem()
			if err := dataToReflectPointer(v, tt.data); err != nil {
				t.Fatalf("dataToReflectPointer(%T, %v) returned error: %v", tt.target, tt.data, err)
			}
			if !reflect.DeepEqual(v.Interface(), tt.want) {
				t.Errorf("dataToReflectPointer() = %#v, want %#v", v.Interface(), tt.want)
			}
		})
	}
}

func TestDataToFloatFieldFromFirestoreInteger(t *testing.T) {
	type measurement struct {
		Value float64
	}
	doc := FirestoreDocument{Fields: map[string]any{
		"Value": map[string]any{"integerValue": "42"},
	}}

	var got measurement
	if err := doc.DataTo(&got); err != nil {
		t.Fatalf("DataTo() returned error: %v", err)
	}
	if got.Value != 42 {
		t.Errorf("DataTo() = %+v, want Value 42", got)
	}
}
