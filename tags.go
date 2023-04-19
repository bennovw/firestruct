// Copyright 2017 Google LLC
// Copyright 2023 Benno Van Waeyenberg
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//	Fork of cloud.google.com/go/firestore with the following changes:
// 	* Removed dependency on the Firestore serverTimestamp

package firestruct

import (
	"fmt"
	"reflect"

	"github.com/bennovw/firestruct/internal/fields"
)

type tagOptions struct {
	omitEmpty bool // do not marshal value if empty
}

// parseTag interprets firestore struct field tags.
func parseTag(t reflect.StructTag) (name string, keep bool, other interface{}, err error) {
	name, keep, opts, err := fields.ParseStandardTag("firestore", t)
	if err != nil {
		return "", false, nil, fmt.Errorf("firestruct: %w", err)
	}
	tagOpts := tagOptions{}
	for _, opt := range opts {
		switch opt {
		case "omitempty":
			tagOpts.omitEmpty = true
		default:
			return "", false, nil, fmt.Errorf("firestruct: unknown tag option: %q", opt)
		}
	}
	return name, keep, tagOpts, nil
}
