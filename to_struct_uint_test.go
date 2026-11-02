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
