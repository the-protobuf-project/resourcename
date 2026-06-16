package main

import (
	"fmt"
	"log"

	"github.com/oklog/ulid/v2"
	"github.com/the-protobuf-project/resourcename"
)

// Artist represents a basic resource
type Artist struct {
	_    struct{} `resource:"//music.example.com/artists/{id}/{name}"`
	ID   string   `resource:"id"`
	Name string   `resource:"name"`
}

// Album represents a nested struct
type Album struct {
	Title string `resource:"title"`
	Year  string `resource:"year"`
}

// ArtistWithAlbum demonstrates nested struct support
type ArtistWithAlbum struct {
	_     struct{} `resource:"//music.example.com/artists/{id}/{album.title}/{album.year}"`
	ID    string   `resource:"id"`
	Album Album    `resource:"album"`
}

func main() {
	fmt.Println("=== Resource Name Marshaling Demo ===")

	fmt.Println("\n1. Basic Example:")
	demoBasic()

	fmt.Println("\n2. Nested Struct Example:")
	demoNested()
}

func demoBasic() {
	ulid := ulid.Make()
	a := &Artist{ID: ulid.String(), Name: "Radiohead"}

	rn, err := resourcename.MarshalResource(a)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Printf("   Marshaled: %s\n", rn)

	a2 := &Artist{}
	err = resourcename.UnmarshalResource(rn, a2)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("   Unmarshaled: ID=%s, Name=%s\n", a2.ID, a2.Name)
}

func demoNested() {
	ulid := ulid.Make()
	a := &ArtistWithAlbum{
		ID:    ulid.String(),
		Album: Album{Title: "In-Rainbows", Year: "2007"},
	}

	rn, err := resourcename.MarshalResource(a)
	if err != nil {
		log.Fatalf("Marshal error: %v", err)
	}
	fmt.Printf("   Marshaled: %s\n", rn)

	a2 := &ArtistWithAlbum{}
	err = resourcename.UnmarshalResource(rn, a2)
	if err != nil {
		log.Fatalf("Unmarshal error: %v", err)
	}
	fmt.Printf("   Unmarshaled: ID=%s, Title=%s, Year=%s\n", a2.ID, a2.Album.Title, a2.Album.Year)
}
