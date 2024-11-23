/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package protosanitizer supports logging of gRPC messages without
// accidentally revealing sensitive fields.
package protosanitizer

import (
	"encoding/json"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// StripSecrets returns a wrapper around the original CSI gRPC message
// which has a Stringer implementation that serializes the message
// as one-line JSON, but without including secret information.
// Instead of the secret value(s), the string "***stripped***" is
// included in the result.
//
// StripSecrets relies on an extension in CSI 1.0 and thus can only
// be used for messages based on that or a more recent spec!
//
// StripSecrets itself is fast and therefore it is cheap to pass the
// result to logging functions which may or may not end up serializing
// the parameter depending on the current log level.
func StripSecrets(msg interface{}) fmt.Stringer {
	return &stripSecrets{msg}
}

type stripSecrets struct {
	msg any
}

func (s *stripSecrets) String() string {
	// First convert to a generic representation. That's less efficient
	// than using reflect directly, but easier to work with.
	var parsed interface{}
	b, err := json.Marshal(s.msg)
	if err != nil {
		return fmt.Sprintf("<<json.Marshal %T: %s>>", s.msg, err)
	}
	msg, ok := s.msg.(proto.Message)
	if !ok {
		return string(b)
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return fmt.Sprintf("<<json.Unmarshal %T: %s>>", s.msg, err)
	}

	// Now remove secrets from the generic representation of the message.
	s.strip(parsed, msg.ProtoReflect())

	// Re-encoded the stripped representation and return that.
	b, err = json.Marshal(parsed)
	if err != nil {
		return fmt.Sprintf("<<json.Marshal %T: %s>>", s.msg, err)
	}
	return string(b)
}

func (s *stripSecrets) strip(parsed interface{}, msg protoreflect.Message) {
	// The corresponding map in the parsed JSON representation.
	parsedFields, ok := parsed.(map[string]interface{})
	if !ok {
		// Probably nil.
		return
	}

	// Walk through all fields and replace those with ***stripped*** that
	// are marked as secret. This relies on protobuf adding "json:" tags
	// on each field where the name matches the field name in the protobuf
	// spec (like volume_capabilities). The field.GetJsonName() method returns
	// a different name (volumeCapabilities) which we don't use.
	msg.Range(func(field protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		name := field.TextName()
		if isCSI1Secret(field) {
			// Overwrite only if already set.
			if _, ok := parsedFields[name]; ok {
				parsedFields[name] = "***stripped***"
			}
		} else if field.Kind() == protoreflect.MessageKind && !field.IsMap() {
			entry := parsedFields[name]
			if field.Cardinality() == protoreflect.Repeated {
				l := v.List()
				// Array of values, like VolumeCapabilities in CreateVolumeRequest.
				for i, entry := range entry.([]interface{}) {
					s.strip(entry, l.Get(i).Message())
				}
			} else {
				// Single value.
				s.strip(entry, v.Message())
			}
		}
		return true
	})
}

// isCSI1Secret uses the csi.E_CsiSecret extension from CSI 1.0 to
// determine whether a field contains secrets.
func isCSI1Secret(desc protoreflect.FieldDescriptor) bool {
	ex := proto.GetExtension(desc.Options(), csi.E_CsiSecret)
	return ex.(bool)
}
