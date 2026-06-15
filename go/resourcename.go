// Package resourcename provides utilities for resource naming following Google Cloud patterns.
//
// This package enables type-safe conversion between Go structs and hierarchical resource
// name strings, supporting Google Cloud's resource naming conventions.
//
// # Resource Marshaling
//
// Convert between Go structs and resource name strings using struct tags. Define a resource
// template on a struct field (typically an anonymous zero-sized field) and map struct fields
// to template placeholders:
//
//	type User struct {
//	    _ struct{} `resource:"//example.com/users/{id}"`
//	    ID string `resource:"id"`
//	}
//
//	u := &User{ID: "u42"}
//	rn, _ := MarshalResource(u)  // "//example.com/users/u42"
//
//	u2 := &User{}
//	UnmarshalResource(rn, u2)    // u2.ID == "u42"
//
// # Supported Types
//
// The following Go types are supported for marshaling and unmarshaling:
//   - string
//   - int, int8, int16, int32, int64
//   - uint, uint8, uint16, uint32, uint64, uintptr
//   - bool
//
// # Multiple Placeholders
//
// Templates can contain multiple placeholders for hierarchical resource names:
//
//	type Device struct {
//	    _ struct{} `resource:"//iot.com/devices/{device_id}/sensors/{sensor_id}"`
//	    DeviceID string `resource:"device_id"`
//	    SensorID string `resource:"sensor_id"`
//	}
//
//	d := &Device{DeviceID: "dev001", SensorID: "temp01"}
//	rn, _ := MarshalResource(d)
//	// "//iot.com/devices/dev001/sensors/temp01"
//
// # Nested Structs
//
// Nested structs are supported using dot notation in template placeholders:
//
//	type Address struct {
//	    City string `resource:"city"`
//	    Zip  string `resource:"zip"`
//	}
//
//	type User struct {
//	    _ struct{} `resource:"//example.com/users/{id}/{address.city}/{address.zip}"`
//	    ID      string  `resource:"id"`
//	    Addr Address `resource:"address"`
//	}
//
//	u := &User{
//	    ID: "u42",
//	    Addr: Address{City: "NYC", Zip: "10001"},
//	}
//	rn, _ := MarshalResource(u)
//	// "//example.com/users/u42/NYC/10001"
//
// # Template Format
//
// Resource templates follow the pattern: //domain.com/path/{placeholder}
//   - Must start with "//"
//   - Placeholders are enclosed in curly braces: {name}
//   - Placeholder values cannot contain forward slashes
//   - Each placeholder must have a corresponding struct field with matching tag
//
// # Error Handling
//
// MarshalResource returns an error if:
//   - The input is nil or not a struct
//   - No resource template is found on the struct
//   - A required field value is missing or nil
//   - A field type is unsupported
//
// UnmarshalResource returns an error if:
//   - The target is nil or not a pointer to struct
//   - No resource template is found on the struct
//   - The resource string doesn't match the template pattern
//   - Type conversion fails (e.g., invalid integer format)
package resourcename
