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
//	type Artist struct {
//	    _ struct{} `resource:"//music.example.com/artists/{id}"`
//	    ID string `resource:"id"`
//	}
//
//	a := &Artist{ID: "radiohead"}
//	rn, _ := MarshalResource(a)  // "//music.example.com/artists/radiohead"
//
//	a2 := &Artist{}
//	UnmarshalResource(rn, a2)    // a2.ID == "radiohead"
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
//	type Track struct {
//	    _ struct{} `resource:"//music.example.com/albums/{album_id}/tracks/{track_id}"`
//	    AlbumID string `resource:"album_id"`
//	    TrackID string `resource:"track_id"`
//	}
//
//	t := &Track{AlbumID: "in-rainbows", TrackID: "15-step"}
//	rn, _ := MarshalResource(t)
//	// "//music.example.com/albums/in-rainbows/tracks/15-step"
//
// # Nested Structs
//
// Nested structs are supported using dot notation in template placeholders:
//
//	type Album struct {
//	    Title string `resource:"title"`
//	    Year  string `resource:"year"`
//	}
//
//	type Artist struct {
//	    _ struct{} `resource:"//music.example.com/artists/{id}/{album.title}/{album.year}"`
//	    ID    string `resource:"id"`
//	    Album Album  `resource:"album"`
//	}
//
//	a := &Artist{
//	    ID: "radiohead",
//	    Album: Album{Title: "In-Rainbows", Year: "2007"},
//	}
//	rn, _ := MarshalResource(a)
//	// "//music.example.com/artists/radiohead/In-Rainbows/2007"
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
