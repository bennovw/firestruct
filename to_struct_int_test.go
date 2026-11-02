package firestruct

import (
	"math"
	"reflect"
	"testing"
)

func TestDataToReflectPointerSignedFromUnsigned(t *testing.T) {
	tests := []struct {
		name    string
		target  any
		data    any
		want    any
		wantErr bool
	}{
		{name: "int from uint", target: new(int), data: uint(1), want: 1},
		{name: "int from uint8", target: new(int), data: uint8(8), want: 8},
		{name: "int16 from uint16", target: new(int16), data: uint16(16), want: int16(16)},
		{name: "int32 from uint32", target: new(int32), data: uint32(32), want: int32(32)},
		{name: "int64 from max int64 as uint64", target: new(int64), data: uint64(math.MaxInt64), want: int64(math.MaxInt64)},
		{name: "int64 overflow from uint64", target: new(int64), data: uint64(math.MaxInt64 + 1), wantErr: true},
		{name: "int8 overflow from uint8", target: new(int8), data: uint8(200), wantErr: true},
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
