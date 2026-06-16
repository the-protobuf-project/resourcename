// Package resourcename provides type-safe marshaling between Go structs and
// hierarchical resource-name strings (Google AIP-122 style).
//
// This is the public entry point for the module. It re-exports the API from the
// implementation package so callers can use the clean import path:
//
//	import "github.com/the-protobuf-project/resourcename"
//
//	type Artist struct {
//	    _    struct{} `resource:"//music.example.com/artists/{id}/{name}"`
//	    ID   string   `resource:"id"`
//	    Name string   `resource:"name"`
//	}
//
//	rn, _ := resourcename.MarshalResource(&Artist{ID: "ar-42", Name: "Radiohead"})
//	// "//music.example.com/artists/ar-42/Radiohead"
//
// The implementation lives in the github.com/the-protobuf-project/resourcename/go
// package; see its documentation for the full template format and supported types.
package resourcename

import impl "github.com/the-protobuf-project/resourcename/go"

// Template is the compiled resource-name template type.
type Template = impl.Template

// MarshalResource converts a struct (or pointer to struct) into a resource string
// using the resource template declared on the struct.
func MarshalResource(v any) (string, error) {
	return impl.MarshalResource(v)
}

// UnmarshalResource parses a resource string and sets fields on v, which must be
// a non-nil pointer to a struct with a resource template.
func UnmarshalResource(resource string, v any) error {
	return impl.UnmarshalResource(resource, v)
}
