package plex_test

import (
	"encoding/json"
	"testing"

	"thde.io/plex"
)

func TestBoolUnmarshalJSON(t *testing.T) {
	tests := []struct {
		in   string
		want plex.Bool
	}{
		{`true`, true},
		{`false`, false},
		{`1`, true},
		{`0`, false},
		{`"true"`, true},
		{`"false"`, false},
		{`"1"`, true},
		{`"0"`, false},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			b := !tt.want
			if err := json.Unmarshal([]byte(tt.in), &b); err != nil {
				t.Fatal(err)
			}
			if b != tt.want {
				t.Errorf("got %v, want %v", b, tt.want)
			}
		})
	}
}

func TestBoolUnmarshalJSONNull(t *testing.T) {
	b := plex.Bool(true)
	if err := json.Unmarshal([]byte(`null`), &b); err != nil {
		t.Fatal(err)
	}
	if !b {
		t.Error("null changed the value")
	}
}

func TestBoolUnmarshalJSONInvalid(t *testing.T) {
	for _, in := range []string{`2`, `-1`, `1.0`, `""`, `"yes"`, `{}`, `[]`} {
		t.Run(in, func(t *testing.T) {
			var b plex.Bool
			if err := json.Unmarshal([]byte(in), &b); err == nil {
				t.Errorf("decoded %s as %v, want error", in, b)
			}
		})
	}
}

func TestBoolMarshalJSON(t *testing.T) {
	got, err := json.Marshal(struct {
		B plex.Bool `json:"b"`
	}{true})
	if err != nil {
		t.Fatal(err)
	}
	if want := `{"b":true}`; string(got) != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

// TestNumericBoolFields decodes fields that Plex Media Server sends as 0 or
// 1 although the specification declares them boolean.
func TestNumericBoolFields(t *testing.T) {
	var section plex.LibrarySection
	if err := json.Unmarshal([]byte(`{"key":"1","type":"movie","hidden":0,"allowSync":true}`), &section); err != nil {
		t.Fatalf("LibrarySection: %v", err)
	}
	if section.Hidden {
		t.Error("LibrarySection.Hidden = true, want false")
	}

	var media plex.Media
	if err := json.Unmarshal([]byte(`{"id":1,"optimizedForStreaming":1,"Part":[{"id":2,"optimizedForStreaming":true}]}`), &media); err != nil {
		t.Fatalf("Media: %v", err)
	}
	if !media.OptimizedForStreaming {
		t.Error("Media.OptimizedForStreaming = false, want true")
	}
}
