package resourcename

import "testing"

type Album struct {
	Title string `resource:"title"`
	Year  string `resource:"year"`
}

type ArtistWithAlbum struct {
	_     struct{} `resource:"//music.example.com/artists/{id}/{album.title}/{album.year}"`
	ID    string   `resource:"id"`
	Album Album    `resource:"album"`
}

type Studio struct {
	Country string `resource:"country"`
	Region  string `resource:"region"`
}

type Label struct {
	Name   string `resource:"name"`
	Studio Studio `resource:"studio"`
}

type Release struct {
	_     struct{} `resource:"//music.example.com/releases/{id}/{label.name}/{label.studio.country}"`
	ID    string   `resource:"id"`
	Label Label    `resource:"label"`
}

func TestNestedStructMarshal(t *testing.T) {
	a := &ArtistWithAlbum{ID: "ar-42", Album: Album{Title: "In-Rainbows", Year: "2007"}}
	rn, err := MarshalResource(a)
	if err != nil {
		t.Fatalf("MarshalResource() error = %v", err)
	}
	expected := "//music.example.com/artists/ar-42/In-Rainbows/2007"
	if rn != expected {
		t.Errorf("got %v, want %v", rn, expected)
	}

	r := &Release{ID: "rel-100", Label: Label{Name: "XL-Recordings", Studio: Studio{Country: "UK", Region: "London"}}}
	rn, err = MarshalResource(r)
	if err != nil {
		t.Fatalf("MarshalResource() error = %v", err)
	}
	expected = "//music.example.com/releases/rel-100/XL-Recordings/UK"
	if rn != expected {
		t.Errorf("got %v, want %v", rn, expected)
	}
}

func TestNestedStructUnmarshal(t *testing.T) {
	a := &ArtistWithAlbum{}
	err := UnmarshalResource("//music.example.com/artists/ar-42/In-Rainbows/2007", a)
	if err != nil {
		t.Fatalf("UnmarshalResource() error = %v", err)
	}
	if a.ID != "ar-42" || a.Album.Title != "In-Rainbows" || a.Album.Year != "2007" {
		t.Errorf("got ID=%v, Title=%v, Year=%v", a.ID, a.Album.Title, a.Album.Year)
	}

	r := &Release{}
	err = UnmarshalResource("//music.example.com/releases/rel-100/XL-Recordings/UK", r)
	if err != nil {
		t.Fatalf("UnmarshalResource() error = %v", err)
	}
	if r.ID != "rel-100" || r.Label.Name != "XL-Recordings" || r.Label.Studio.Country != "UK" {
		t.Errorf("got ID=%v, Name=%v, Country=%v", r.ID, r.Label.Name, r.Label.Studio.Country)
	}
}

func TestNestedStructRoundTrip(t *testing.T) {
	a := &ArtistWithAlbum{ID: "ar-42", Album: Album{Title: "OK-Computer", Year: "1997"}}
	rn, _ := MarshalResource(a)
	a2 := &ArtistWithAlbum{}
	if err := UnmarshalResource(rn, a2); err != nil {
		t.Errorf("unmarshal failed: %v", err)
	}
	rn2, _ := MarshalResource(a2)
	if rn != rn2 {
		t.Errorf("round trip failed: %v != %v", rn, rn2)
	}

	r := &Release{ID: "rel-200", Label: Label{Name: "Parlophone", Studio: Studio{Country: "UK", Region: "E"}}}
	rn, _ = MarshalResource(r)
	r2 := &Release{}
	if err := UnmarshalResource(rn, r2); err != nil {
		t.Errorf("unmarshal failed: %v", err)
	}
	rn2, _ = MarshalResource(r2)
	if rn != rn2 {
		t.Errorf("round trip failed: %v != %v", rn, rn2)
	}
}
