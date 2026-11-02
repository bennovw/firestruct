package firestruct

import (
	"reflect"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := []struct {
		name     string
		tag      reflect.StructTag
		wantName string
		wantKeep bool
		wantOpts tagOptions
		wantErr  bool
	}{
		{name: "no tag", tag: ``, wantKeep: true},
		{name: "name only", tag: `firestore:"title"`, wantName: "title", wantKeep: true},
		{name: "omitempty", tag: `firestore:"title,omitempty"`, wantName: "title", wantKeep: true, wantOpts: tagOptions{omitEmpty: true}},
		{name: "ignored field", tag: `firestore:"-"`},
		{name: "ignored field with options", tag: `firestore:"-,omitempty"`, wantErr: true},
		{name: "unknown option", tag: `firestore:"title,bogus"`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, keep, other, err := parseTag(tt.tag)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseTag(%q) error = %v, wantErr %v", tt.tag, err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if name != tt.wantName || keep != tt.wantKeep {
				t.Errorf("parseTag(%q) = %q, %v; want %q, %v", tt.tag, name, keep, tt.wantName, tt.wantKeep)
			}
			if other != nil && other.(tagOptions) != tt.wantOpts {
				t.Errorf("parseTag(%q) options = %+v, want %+v", tt.tag, other, tt.wantOpts)
			}
		})
	}
}
