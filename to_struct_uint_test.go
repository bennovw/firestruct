package firestruct

import (
	"reflect"
	"testing"
)

func TestDataToReflectPointerUnsigned(t *testing.T) {
	tests := []struct {
		name    string
		target  any
		data    any
		want    any
		wantErr bool
	}{
		{name: "uint from uint", target: new(uint), data: uint(7), want: uint(7)},
		{name: "uint64 from uint64", target: new(uint64), data: uint64(1 << 63), want: uint64(1 << 63)},
		{name: "uint8 from uint", target: new(uint8), data: uint(8), want: uint8(8)},
		{name: "uint8 overflow from uint", target: new(uint8), data: uint(300), wantErr: true},
		{name: "uint8 from int", target: new(uint8), data: 8, want: uint8(8)},
		{name: "uint32 from int64", target: new(uint32), data: int64(32), want: uint32(32)},
		{name: "uint64 from int", target: new(uint64), data: 64, want: uint64(64)},
		{name: "uint16 overflow from int", target: new(uint16), data: 70000, wantErr: true},
		{name: "uint from negative int", target: new(uint), data: -1, wantErr: true},
		{name: "uint8 from string", target: new(uint8), data: "8", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.target).Elem()
			err := dataToReflectPointer(v, tt.data)
			if (err != nil) != tt.wantErr {
				t.Fatalf("dataToReflectPointer(%T, %v) error = %v, wantErr %v", tt.target, tt.data, err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(v.Interface(), tt.want) {
				t.Errorf("dataToReflectPointer() = %#v, want %#v", v.Interface(), tt.want)
			}
		})
	}
}

func TestDataToUnsignedFieldsFromFirestore(t *testing.T) {
	type counters struct {
		Small uint8
		Large uint64
	}
	doc := FirestoreDocument{Fields: map[string]any{
		"Small": map[string]any{"integerValue": "8"},
		"Large": map[string]any{"integerValue": "9000000000"},
	}}

	var got counters
	if err := doc.DataTo(&got); err != nil {
		t.Fatalf("DataTo() returned error: %v", err)
	}
	if want := (counters{Small: 8, Large: 9000000000}); got != want {
		t.Errorf("DataTo() = %+v, want %+v", got, want)
	}
}
