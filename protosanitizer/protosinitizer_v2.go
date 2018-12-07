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
	"fmt"
	"reflect"
	"strings"

	csipb "github.com/container-storage-interface/spec/lib/go/csi"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	descr "github.com/golang/protobuf/protoc-gen-go/descriptor"
)

var (
	sanitizeFields  []descr.FieldDescriptorProto
	sanitizeMessage = "***Sanitized***"
)

// SanitizeMsg scans the structure for fields tagged with csi_secure and removes their values.
func SanitizeMsg(msg interface{}, additionalFieldsToSanitize ...string) string {

	if err := buildSanitizationList(msg, additionalFieldsToSanitize); err != nil {
		return fmt.Sprintf("failed to build sanitization list with error: %+v\n", err)
	}

	if len(sanitizeFields) == 0 {
		return fmt.Sprintf("%+v", msg)
	}

	copy := reflect.New(reflect.TypeOf(msg))
	traverseAndSanitize(reflect.ValueOf(msg), copy.Elem(), false)

	return fmt.Sprintf("%+v", copy.Elem())
}

func traverseAndSanitize(obj, copy reflect.Value, sanitize bool) {
	switch obj.Kind() {
	case reflect.Ptr:
		pv := reflect.Indirect(obj)
		if !pv.IsValid() {
			return
		}
		copy.Set(reflect.New(pv.Type()))
		traverseAndSanitize(pv, reflect.Indirect(copy), sanitize)
	case reflect.Interface:
		objV := obj.Elem()
		if !objV.IsValid() {
			return
		}
		copyV := reflect.New(objV.Type()).Elem()
		traverseAndSanitize(objV, copyV, sanitize)
		copy.Set(copyV)
	case reflect.Struct:
		for i := 0; i < obj.NumField(); i++ {
			if !needSanitize(obj.Type().Field(i).Name) {
				sanitize = false
				traverseAndSanitize(obj.Field(i), copy.Field(i), sanitize)
			} else {
				sanitize = true
				traverseAndSanitize(obj.Field(i), copy.Field(i), sanitize)
			}
		}
	case reflect.Slice:
		copy.Set(reflect.MakeSlice(obj.Type(), obj.Len(), obj.Cap()))
		for i := 0; i < obj.Len(); i++ {
			traverseAndSanitize(obj.Index(i), copy.Index(i), sanitize)
		}
	case reflect.Map:
		copy.Set(reflect.MakeMap(obj.Type()))
		for _, key := range obj.MapKeys() {
			objV := obj.MapIndex(key)
			copyV := reflect.New(objV.Type()).Elem()
			traverseAndSanitize(objV, copyV, sanitize)
			copy.SetMapIndex(key, copyV)
		}
	case reflect.String:
		if sanitize {
			copy.SetString(sanitizeMessage)
		} else {
			copy.Set(obj)
		}
	case reflect.Uint8:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int:
		fallthrough
	case reflect.Int64:
		if sanitize {
			copy.SetInt(int64(0))
		} else {
			copy.Set(obj)
		}
	default:
		copy.Set(obj)
	}
}

func buildSanitizationList(pb interface{}, additionalFieldsToSanitize []string) error {
	if _, ok := pb.(descriptor.Message); !ok {
		return fmt.Errorf("invalid message")
	}
	_, md := descriptor.ForMessage(pb.(descriptor.Message))
	fields := md.GetField()
	if fields == nil {
		return fmt.Errorf("fail to get a list of fields")
	}
	for _, field := range fields {
		opt, err := proto.GetExtension(field.Options, csipb.E_CsiSecret)
		if err == nil {
			_, ok := opt.(*bool)
			if ok {
				sanitizeFields = append(sanitizeFields, *field)
			}
		}
		if isOnAdditionalFieldsList(*field.Name, additionalFieldsToSanitize) {
			sanitizeFields = append(sanitizeFields, *field)
		}
	}
	return nil
}

func isOnAdditionalFieldsList(field string, additionalFieldsToSanitize []string) bool {
	n := strings.ToLower(field)
	for _, f := range additionalFieldsToSanitize {
		if f == n {
			return true
		}
	}
	return false
}

func needSanitize(n string) bool {
	name := strings.ToLower(n)
	for _, f := range sanitizeFields {
		if name == *f.Name {
			return true
		}
	}
	return false
}
