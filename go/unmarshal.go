// unmarshal.go
package resourcename

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// UnmarshalResource parses a resource string and sets fields on v (pointer to struct).
// The struct must have a resource template defined either via:
// 1. A ResourceTemplate() string method
// 2. A struct field with tag `resource:"//domain.com/Type/{placeholder}"`
//
// Fields to be unmarshaled must have `resource:"placeholder"` tags matching the template.
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
//	u := &User{}
//	err := resourcename.UnmarshalResource("//protobuf_project.com/User/u42/Ria/29", u)
//	// u.ID == "u42", u.Name == "Ria", u.Age == 29
func UnmarshalResource(resource string, v interface{}) error {
	if v == nil {
		return fmt.Errorf("nil target")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("UnmarshalResource requires a non-nil pointer to a struct")
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return fmt.Errorf("UnmarshalResource requires pointer to struct")
	}
	rt := rv.Type()

	tmplStr, err := getTemplateForType(rv)
	if err != nil {
		return err
	}
	if tmplStr == "" {
		return fmt.Errorf("no resource template found on type %s", rt.String())
	}
	tpl, err := build(tmplStr)
	if err != nil {
		return err
	}

	parsed, err := tpl.Parse(resource)
	if err != nil {
		return err
	}

	return setFieldValues(rv, rt, "", parsed)
}

// setFieldValues recursively sets field values on a struct, supporting nested structs.
// prefix is used for dot notation (e.g., "address.city")
func setFieldValues(rv reflect.Value, rt reflect.Type, prefix string, parsed map[string]string) error {
	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)
		tag, ok := sf.Tag.Lookup("resource")
		if !ok {
			continue
		}
		if strings.HasPrefix(tag, "//") || (strings.Contains(tag, "{") && strings.Contains(tag, "}")) {
			continue
		}
		// validate that tag is not empty
		if tag == "" {
			return fmt.Errorf("field %s has empty resource tag", sf.Name)
		}

		fv := rv.Field(i)
		if !fv.CanSet() {
			continue
		}

		// follow pointers
		if fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			fv = fv.Elem()
		}

		// Handle nested structs
		if fv.Kind() == reflect.Struct {
			nestedPrefix := tag
			if prefix != "" {
				nestedPrefix = prefix + "." + tag
			}
			if err := setFieldValues(fv, fv.Type(), nestedPrefix, parsed); err != nil {
				return fmt.Errorf("field %s: %w", sf.Name, err)
			}
			continue
		}

		// Look up value with prefix
		key := tag
		if prefix != "" {
			key = prefix + "." + tag
		}
		valStr, found := parsed[key]
		if !found {
			continue
		}

		if err := setValueFromString(fv, valStr); err != nil {
			return fmt.Errorf("field %s: %w", sf.Name, err)
		}
	}
	return nil
}

// setValueFromString sets a reflect.Value from its string representation.
func setValueFromString(fv reflect.Value, s string) error {
	switch fv.Kind() {
	case reflect.String:
		fv.SetString(s)
		return nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetInt(n)
		return nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		fv.SetUint(n)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		fv.SetBool(b)
		return nil
	default:
		return fmt.Errorf("unsupported kind for unmarshaling: %s", fv.Kind())
	}
}
