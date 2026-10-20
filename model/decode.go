package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// Decode checks the complete JSON value, rejects unknown struct fields, null
// for non-nullable fields, and missing fields marked required in the contract.
// Errors contain schema locations only, never values from private payloads.
func Decode(raw []byte, target any) error {
	if target == nil || reflect.TypeOf(target).Kind() != reflect.Pointer || reflect.ValueOf(target).IsNil() {
		return fmt.Errorf("decode target must be a non-nil pointer")
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	dec.DisallowUnknownFields()
	if err := dec.Decode(target); err != nil {
		return fmt.Errorf("response does not match %T", target)
	}
	if err := dec.Decode(new(any)); err != io.EOF {
		return fmt.Errorf("expected one JSON value")
	}
	return validateShape(raw, reflect.TypeOf(target).Elem(), "data")
}

var unmarshalerType = reflect.TypeFor[json.Unmarshaler]()

func validateShape(raw []byte, typ reflect.Type, path string) error {
	raw = bytes.TrimSpace(raw)
	if typ.Kind() == reflect.Pointer {
		if bytes.Equal(raw, []byte("null")) {
			return nil
		}
		return validateShape(raw, typ.Elem(), path)
	}
	// Custom contracts validate their own shape (tuples and compatibility unions).
	if reflect.PointerTo(typ).Implements(unmarshalerType) {
		return nil
	}
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return fmt.Errorf("%s must not be null", path)
	}
	switch typ.Kind() {
	case reflect.Struct:
		var obj map[string]json.RawMessage
		if json.Unmarshal(raw, &obj) != nil || obj == nil {
			return fmt.Errorf("%s must be an object", path)
		}
		known := map[string]bool{}
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if !f.IsExported() {
				continue
			}
			name := strings.Split(f.Tag.Get("json"), ",")[0]
			if name == "" {
				name = f.Name
			}
			known[name] = true
		}
		for name := range obj {
			if !known[name] {
				return fmt.Errorf("%s contains an unknown field", path)
			}
		}
		for i := 0; i < typ.NumField(); i++ {
			f := typ.Field(i)
			if !f.IsExported() {
				continue
			}
			name := strings.Split(f.Tag.Get("json"), ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				name = f.Name
			}
			v, ok := obj[name]
			if !ok {
				if f.Tag.Get("required") == "true" {
					return fmt.Errorf("%s.%s is required", path, name)
				}
				continue
			}
			if f.Tag.Get("required") == "true" && bytes.Equal(bytes.TrimSpace(v), []byte("null")) {
				return fmt.Errorf("%s.%s must not be null", path, name)
			}
			if err := validateShape(v, f.Type, path+"."+name); err != nil {
				return err
			}
		}
	case reflect.Slice, reflect.Array:
		var arr []json.RawMessage
		if json.Unmarshal(raw, &arr) != nil || arr == nil {
			return fmt.Errorf("%s must be an array", path)
		}
		for i, v := range arr {
			if err := validateShape(v, typ.Elem(), fmt.Sprintf("%s[%d]", path, i)); err != nil {
				return err
			}
		}
	case reflect.Map:
		var obj map[string]json.RawMessage
		if json.Unmarshal(raw, &obj) != nil || obj == nil {
			return fmt.Errorf("%s must be an object", path)
		}
		for _, v := range obj {
			if err := validateShape(v, typ.Elem(), path+".*"); err != nil {
				return err
			}
		}
	}
	return nil
}

// NoData is the acknowledgement of a mutation, never a fill or a snapshot.
// These pages omit data, return null, or show an empty object.
type NoData struct{}

func (*NoData) UnmarshalJSON(raw []byte) error {
	if string(raw) == "null" || string(raw) == "{}" {
		return nil
	}
	return fmt.Errorf("expected an empty acknowledgement")
}
