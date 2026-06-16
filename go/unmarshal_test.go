package resourcename

import (
	"testing"
)

func TestUnmarshalResource(t *testing.T) {
	tests := []struct {
		name     string
		resource string
		target   interface{}
		check    func(t *testing.T, v interface{})
		wantErr  bool
	}{
		{
			name:     "basic artist unmarshal",
			resource: "//music.example.com/artists/ar-42/Radiohead/1985",
			target:   &Artist{},
			check: func(t *testing.T, v interface{}) {
				a := v.(*Artist)
				if a.ID != "ar-42" {
					t.Errorf("ID = %v, want ar-42", a.ID)
				}
				if a.Name != "Radiohead" {
					t.Errorf("Name = %v, want Radiohead", a.Name)
				}
				if a.Year != 1985 {
					t.Errorf("Year = %v, want 1985", a.Year)
				}
			},
			wantErr: false,
		},
		{
			name:     "recording unmarshal",
			resource: "//music.example.com/genres/rock/recordings/OK-Computer",
			target:   &Recording{},
			check: func(t *testing.T, v interface{}) {
				r := v.(*Recording)
				if r.Genre != "rock" {
					t.Errorf("Genre = %v, want rock", r.Genre)
				}
				if r.Title != "OK-Computer" {
					t.Errorf("Title = %v, want OK-Computer", r.Title)
				}
			},
			wantErr: false,
		},
		{
			name:     "track unmarshal",
			resource: "//music.example.com/albums/in-rainbows/tracks/15-step",
			target:   &Track{},
			check: func(t *testing.T, v interface{}) {
				tr := v.(*Track)
				if tr.AlbumID != "in-rainbows" {
					t.Errorf("AlbumID = %v, want in-rainbows", tr.AlbumID)
				}
				if tr.TrackID != "15-step" {
					t.Errorf("TrackID = %v, want 15-step", tr.TrackID)
				}
			},
			wantErr: false,
		},
		{
			name:     "invalid resource format",
			resource: "invalid-resource",
			target:   &Artist{},
			check:    nil,
			wantErr:  true,
		},
		{
			name:     "nil target",
			resource: "//music.example.com/artists/ar-42/Radiohead/1985",
			target:   nil,
			check:    nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UnmarshalResource(tt.resource, tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, tt.target)
			}
		})
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
	}{
		{
			name: "artist round trip",
			input: &Artist{
				ID:   "ar-42",
				Name: "Radiohead",
				Year: 1985,
			},
		},
		{
			name: "recording round trip",
			input: &Recording{
				Genre: "jazz",
				Title: "Kind-Of-Blue",
			},
		},
		{
			name: "track round trip",
			input: &Track{
				AlbumID:  "debut",
				TrackID:  "human-behaviour",
				Explicit: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			resource, err := MarshalResource(tt.input)
			if err != nil {
				t.Fatalf("MarshalResource() error = %v", err)
			}

			// Unmarshal
			var target interface{}
			switch tt.input.(type) {
			case *Artist:
				target = &Artist{}
			case *Recording:
				target = &Recording{}
			case *Track:
				target = &Track{}
			}

			err = UnmarshalResource(resource, target)
			if err != nil {
				t.Fatalf("UnmarshalResource() error = %v", err)
			}

			// Marshal again and compare
			resource2, err := MarshalResource(target)
			if err != nil {
				t.Fatalf("MarshalResource() second time error = %v", err)
			}

			if resource != resource2 {
				t.Errorf("Round trip failed: %v != %v", resource, resource2)
			}
		})
	}
}
