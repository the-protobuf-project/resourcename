package resourcename

import (
	"testing"
)

type Artist struct {
	_    struct{} `resource:"//music.example.com/artists/{id}/{name}/{year}"`
	ID   string   `resource:"id"`
	Name string   `resource:"name"`
	Year int      `resource:"year"`
}

type Recording struct {
	_     struct{} `resource:"//music.example.com/genres/{genre}/recordings/{title}"`
	Genre string   `resource:"genre"`
	Title string   `resource:"title"`
}

type Track struct {
	_        struct{} `resource:"//music.example.com/albums/{album_id}/tracks/{track_id}"`
	AlbumID  string   `resource:"album_id"`
	TrackID  string   `resource:"track_id"`
	Explicit bool     `resource:"explicit"`
}

func TestMarshalResource(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
		wantErr  bool
	}{
		{
			name: "basic artist struct",
			input: &Artist{
				ID:   "ar-42",
				Name: "Radiohead",
				Year: 1985,
			},
			expected: "//music.example.com/artists/ar-42/Radiohead/1985",
			wantErr:  false,
		},
		{
			name: "artist struct by value",
			input: Artist{
				ID:   "ar-100",
				Name: "Bjork",
				Year: 1993,
			},
			expected: "//music.example.com/artists/ar-100/Bjork/1993",
			wantErr:  false,
		},
		{
			name: "recording struct",
			input: &Recording{
				Genre: "rock",
				Title: "OK-Computer",
			},
			expected: "//music.example.com/genres/rock/recordings/OK-Computer",
			wantErr:  false,
		},
		{
			name: "track with bool",
			input: &Track{
				AlbumID:  "in-rainbows",
				TrackID:  "15-step",
				Explicit: true,
			},
			expected: "//music.example.com/albums/in-rainbows/tracks/15-step",
			wantErr:  false,
		},
		{
			name:     "nil input",
			input:    nil,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := MarshalResource(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("MarshalResource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && result != tt.expected {
				t.Errorf("MarshalResource() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestMarshalResource_EdgeCases(t *testing.T) {
	t.Run("zero values", func(t *testing.T) {
		a := &Artist{
			ID:   "",
			Name: "",
			Year: 0,
		}
		result, err := MarshalResource(a)
		if err != nil {
			t.Errorf("MarshalResource() error = %v", err)
		}
		expected := "//music.example.com/artists///0"
		if result != expected {
			t.Errorf("MarshalResource() = %v, want %v", result, expected)
		}
	})

	t.Run("special characters in name", func(t *testing.T) {
		a := &Artist{
			ID:   "ar-42",
			Name: "Sigur-Ros",
			Year: 1994,
		}
		result, err := MarshalResource(a)
		if err != nil {
			t.Errorf("MarshalResource() error = %v", err)
		}
		expected := "//music.example.com/artists/ar-42/Sigur-Ros/1994"
		if result != expected {
			t.Errorf("MarshalResource() = %v, want %v", result, expected)
		}
	})
}
