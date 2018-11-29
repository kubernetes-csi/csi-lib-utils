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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"
	protobufdescriptor "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/pkg/errors"
)

// StripSecrets returns a wrapper around the original CSI gRPC message
// which has a Stringer implementation that serializes the message
// similar to gRPC, but without including secret information.
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
	return &stripSecrets{msg, isCSI1Secret}
}

// StripSecretsCSI03 is like StripSecrets, except that it works
// for messages based on CSI 0.3 and older. It does not work
// for CSI 1.0, use StripSecrets for that.
func StripSecretsCSI03(msg interface{}) fmt.Stringer {
	return &stripSecrets{msg, isCSI03Secret}
}

type stripped struct{}

func (s stripped) String() string {
	return "***stripped***"
}

type stripSecrets struct {
	msg interface{}

	isSecretField func(field *protobuf.FieldDescriptorProto) bool
}

func (s *stripSecrets) String() string {
	// First convert to a generic representation. That's less efficient
	// than using reflect directly, but easier to work with.
	var parsed interface{}
	b, err := json.Marshal(s.msg)
	if err != nil {
		return fmt.Sprintf("<<json.Marshal %T: %s>>", s.msg, err)
	}
	if err := json.Unmarshal(b, &parsed); err != nil {
		return fmt.Sprintf("<<json.Unmarshal %T: %s>>", s.msg, err)
	}

	// Now remove secrets from the generic representation of the message.
	if err := s.strip(parsed, s.msg); err != nil {
		return fmt.Sprintf("<<error: %s>>", err)
	}

	// Re-encoded the stripped representation and return that.
	var buf bytes.Buffer
	marshal(&buf, parsed)
	return buf.String()
}

// marshall is mimicking the output of the compact text marshaller
// from https://github.com/golang/protobuf/blob/master/proto/text.go. It
// can't be exactly the same because it works on the result of json.Unmarshal
// into generic types (map, array, int/string/bool).
func marshal(out io.Writer, parsed interface{}) {
	switch parsed := parsed.(type) {
	case map[string]interface{}:
		out.Write([]byte("<"))
		var keys []string
		for k := range parsed {
			keys = append(keys, k)
		}
		// Ensure consistent alphabetic ordering.
		// Ordering by field number would be nicer, but we don't
		// have that information here.
		sort.Strings(keys)
		for i, k := range keys {
			if isFieldName(k) {
				// No quotation marks round simple
				// strings that are likely to be field
				// names. We can't be sure 100%,
				// because both structs and real maps
				// end up as maps after
				// json.Unmarshal.
				out.Write([]byte(k))
			} else {
				out.Write([]byte(fmt.Sprintf("%q", k)))
			}
			out.Write([]byte(":"))
			marshal(out, parsed[k])
			// Avoid redundant space after last element.
			if i+1 < len(keys) {
				out.Write([]byte(" "))
			}
		}
		out.Write([]byte(">"))
	case []interface{}:
		// gRPC uses < for repeating elements. We use
		// [ ] because it is a bit more readable.
		out.Write([]byte("["))
		for i, v := range parsed {
			marshal(out, v)
			// Avoid redundant space after last element.
			if i+1 < len(parsed) {
				out.Write([]byte(" "))
			}
		}
		out.Write([]byte("]"))
	case string:
		fmt.Fprintf(out, "%q", parsed)
	default:
		fmt.Fprint(out, parsed)
	}
}

// isFieldName returns true for strings that start with a-zA-Z and
// are followed by those, digits or underscore.
func isFieldName(str string) bool {
	return fieldNameRe.MatchString(str)
}

var fieldNameRe = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_]*`)

func (s *stripSecrets) strip(parsed interface{}, msg interface{}) error {
	protobufMsg, ok := msg.(descriptor.Message)
	if !ok {
		// Not a protobuf message, nothing to strip.
		return nil
	}

	// The corresponding map in the parsed JSON representation.
	parsedFields, ok := parsed.(map[string]interface{})
	if !ok {
		// Probably nil, nothing to strip.
		return nil
	}

	// Walk through all fields and replace those with ***stripped*** that
	// are marked as secret. This relies on protobuf adding "json:" tags
	// on each field where the name matches the field name in the protobuf
	// spec (like volume_capabilities). The field.GetJsonName() method returns
	// a different name (volumeCapabilities) which we don't use.
	_, md := descriptor.ForMessage(protobufMsg)
	fields := md.GetField()
	if fields != nil {
		for _, field := range fields {
			if s.isSecretField(field) {
				// Overwrite only if already set.
				if _, ok := parsedFields[field.GetName()]; ok {
					parsedFields[field.GetName()] = stripped{}
				}
				continue
			}

			// Not stripped. Decide whether we need to
			// recursively strip the message(s) that the
			// field contains.

			if field.GetType() != protobuf.FieldDescriptorProto_TYPE_MESSAGE {
				// No need to recurse into plain types.
				continue
			}

			// When we get here, the type name is something
			// like ".csi.v1.CapacityRange" (leading dot!)
			// and looking up "csi.v1.CapacityRange"
			// returns the type of a pointer to a pointer
			// to CapacityRange. We need a pointer to such
			// a value for recursive stripping.
			typeName := field.GetTypeName()
			if strings.HasPrefix(typeName, ".") {
				typeName = typeName[1:]
			}
			t := proto.MessageType(typeName)
			// Shouldn't happen, but better check
			// anyway instead of panicking.
			if t == nil {
				return errors.Errorf("%s: unknown type", typeName)
			}
			var v reflect.Value
			switch t.Kind() {
			case reflect.Map:
				ptrType := t.Elem()
				if ptrType.Kind() != reflect.Ptr {
					// map to plain type, nothing to recurse into
					continue
				}
				v = reflect.New(t.Elem().Elem())
			case reflect.Ptr:
				v = reflect.New(t.Elem())
			default:
				return errors.Errorf("%s: has no elements", typeName)
			}
			i := v.Interface()

			if field.OneofIndex != nil {
				// A oneof field doesn't have json tags in the generated .pb.go.
				// Therefore the parsedFields is different and we need to recurse differently:
				// - the entry is named like the Go field (upper first character, no underscores)
				// - it contains a map with the name of the individual candidates, again
				//   with Go naming
				// Example for proto field "volume" inside "VolumeContentSource":
				// map[string]interface {} [
				//    "Type": map[string]interface {} [
				//       "Volume": *(*"interface {}")(0xc4201c0988),
				//    ],
				// ]
				if len(md.OneofDecl) <= int(*field.OneofIndex) {
					// Shouldn't happen, bail out.
					return errors.Errorf("invalid oneof index in %v", field)
				}
				oneof := md.OneofDecl[int(*field.OneofIndex)]
				if oneof.Name == nil {
					return errors.Errorf("invalid oneof for %s, no name: %v", typeName, oneof)
				}
				jsonName := upperFirst(*oneof.Name)
				entry, ok := parsedFields[jsonName]
				if !ok || entry == nil {
					// Not set.
					continue
				}
				oneofMap, ok := entry.(map[string]interface{})
				if !ok {
					return errors.Errorf("unexpected type %T in JSON for %s", entry, typeName)
				}
				if field.JsonName == nil {
					return errors.Errorf("invalid field %v, no name", field)
				}
				entry, ok = oneofMap[upperFirst(*field.JsonName)]
				if !ok {
					// Oneof does not contain this particular field.
					continue
				}
				// Finally recurse into a particular oneof value.
				if err := s.strip(entry, i); err != nil {
					return errors.Wrap(err, typeName)
				}
				continue
			}

			entry := parsedFields[field.GetName()]
			if slice, ok := entry.([]interface{}); ok {
				// Array of values, like VolumeCapabilities in CreateVolumeRequest.
				for _, entry := range slice {
					if err := s.strip(entry, i); err != nil {
						return errors.Wrap(err, typeName)
					}
				}
				continue
			}
			if mapping, ok := entry.(map[string]interface{}); ok {
				// All maps in protobuf are string to something maps in JSON.
				for _, entry := range mapping {
					if err := s.strip(entry, i); err != nil {
						return errors.Wrap(err, typeName)
					}
				}
			}

			// Single value.
			if err := s.strip(entry, i); err != nil {
				return errors.Wrap(err, typeName)
			}
		}
	}
	return nil
}

func upperFirst(str string) string {
	if len(str) >= 2 {
		return strings.ToUpper(str[:1]) + str[1:]
	}
	return strings.ToUpper(str)
}

// isCSI1Secret uses the csi.E_CsiSecret extension from CSI 1.0 to
// determine whether a field contains secrets.
func isCSI1Secret(field *protobuf.FieldDescriptorProto) bool {
	ex, err := proto.GetExtension(field.Options, e_CsiSecret)
	return err == nil && ex != nil && *ex.(*bool)
}

// Copied from the CSI 1.0 spec (https://github.com/container-storage-interface/spec/blob/37e74064635d27c8e33537c863b37ccb1182d4f8/lib/go/csi/csi.pb.go#L4520-L4527)
// to avoid a package dependency that would prevent usage of this package
// in repos using an older version of the spec.
//
// Future revision of the CSI spec must not change this extensions, otherwise
// they will break filtering in binaries based on the 1.0 version of the spec.
var e_CsiSecret = &proto.ExtensionDesc{
	ExtendedType:  (*protobufdescriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         1059,
	Name:          "csi.v1.csi_secret",
	Tag:           "varint,1059,opt,name=csi_secret,json=csiSecret",
	Filename:      "github.com/container-storage-interface/spec/csi.proto",
}

// isCSI03Secret relies on the naming convention in CSI <= 0.3
// to determine whether a field contains secrets.
func isCSI03Secret(field *protobuf.FieldDescriptorProto) bool {
	return strings.HasSuffix(field.GetName(), "_secrets")
}
