// marshal.go
package resourcename

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// MarshalResource converts a struct (or pointer to struct) into a resource string.
// The struct must have a resource template defined either via:
// 1. A ResourceTemplate() string method
// 2. A struct field with tag `resource:"//domain.com/Type/{placeholder}"`
//
// Fields to be marshaled must have `resource:"placeholder"` tags matching the template.
// Nested structs are supported using dot notation (e.g., {address.city}).
//
// Example:
//
//	type User struct {
//	    _ struct{} `resource:"//protobuf_project.com/User/{id}/{name}/{age}"`
//	    ID   string `resource:"id"`
//	    Name string `resource:"name"`
//	    Age  int    `resource:"age"`
//	}
//
//	u := &User{ID: "u42", Name: "Ria", Age: 29}
//	rn, err := resourcename.MarshalResource(u)
//	// rn == "//protobuf_project.com/User/u42/Ria/29"
func MarshalResource(v interface{}) (string, error) {
	if v == nil {
		return "", fmt.Errorf("nil value")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr && !rv.IsNil() {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return "", fmt.Errorf("MarshalResource: expecting struct or pointer to struct")
	}
	rt := rv.Type()

	tmplStr, err := getTemplateForType(rv)
	if err != nil {
		return "", err
	}
	if tmplStr == "" {
		return "", fmt.Errorf("no resource template found on type %s", rt.String())
	}
	tpl, err := build(tmplStr)
	if err != nil {
		return "", err
	}

	values := map[string]string{}
	if err := collectFieldValues(rv, rt, "", values); err != nil {
		return "", err
	}

	return tpl.Generate(values)
}

// collectFieldValues recursively collects field values from a struct, supporting nested structs.
// prefix is used for dot notation (e.g., "address.city")
func collectFieldValues(rv reflect.Value, rt reflect.Type, prefix string, values map[string]string) error {
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		tag, ok := sf.Tag.Lookup("resource")
		if !ok {
			continue
		}
		// skip template carrier field
		if strings.HasPrefix(tag, "//") || (strings.Contains(tag, "{") && strings.Contains(tag, "}")) {
			continue
		}
		// validate that tag is not empty
		if tag == "" {
			return fmt.Errorf("field %s has empty resource tag", sf.Name)
		}
		// only exported fields
		if sf.PkgPath != "" {
			continue
		}
		fv := rv.Field(i)
		// follow pointers
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				return fmt.Errorf("nil pointer for field %s", sf.Name)
			}
			fv = fv.Elem()
		}

		// Handle nested structs
		if fv.Kind() == reflect.Struct {
			// Recursively collect from nested struct
			nestedPrefix := tag
			if prefix != "" {
				nestedPrefix = prefix + "." + tag
			}
			if err := collectFieldValues(fv, fv.Type(), nestedPrefix, values); err != nil {
				return fmt.Errorf("field %s: %w", sf.Name, err)
			}
			continue
		}

		// basic conversion
		s, err := valueToString(fv)
		if err != nil {
			return fmt.Errorf("field %s: %w", sf.Name, err)
		}

		// Use prefix for nested fields
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}
		values[key] = s
	}
	return nil
}

// valueToString converts a reflect.Value to its string representation.
func valueToString(v reflect.Value) (string, error) {
	switch v.Kind() {
	case reflect.String:
		return v.String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10), nil
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), nil
	default:
		return "", fmt.Errorf("unsupported kind for marshaling: %s", v.Kind())
	}
}

// getTemplateForType looks for a template in two places (in order):
// 1) Method ResourceTemplate() string on the type (value or pointer receiver)
// 2) A struct field with tag `resource:"<template>"` (commonly an anonymous zero-sized field)
func getTemplateForType(rv reflect.Value) (string, error) {
	// Ensure we have a value
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv = reflect.New(rv.Type().Elem()).Elem()
		} else {
			rv = rv.Elem()
		}
	}
	rt := rv.Type()

	// 1) method on value
	if m := rv.MethodByName("ResourceTemplate"); m.IsValid() {
		if m.Type().NumIn() == 0 && m.Type().NumOut() == 1 && m.Type().Out(0).Kind() == reflect.String {
			out := m.Call(nil)
			if len(out) == 1 {
				return out[0].String(), nil
			}
		}
	}

	// 1b) method on pointer receiver
	ptr := reflect.New(rt)
	if m := ptr.MethodByName("ResourceTemplate"); m.IsValid() {
		if m.Type().NumIn() == 0 && m.Type().NumOut() == 1 && m.Type().Out(0).Kind() == reflect.String {
			out := m.Call(nil)
			if len(out) == 1 {
				return out[0].String(), nil
			}
		}
	}

	// 2) struct field tag
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		if tag, ok := sf.Tag.Lookup("resource"); ok {
			// consider it a template if it looks like one
			if strings.HasPrefix(tag, "//") || (strings.Contains(tag, "{") && strings.Contains(tag, "}")) {
				return tag, nil
			}
		}
	}
	return "", nil
}
