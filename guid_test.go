package plex_test

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"

	"thde.io/plex"
)

// TestMetadataGUID lists a show section with includeGuids=1, which returns
// the external identifiers of each show as strings. The specification types
// them as tags with integer IDs.
func TestMetadataGUID(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"MediaContainer":{"size":1,"Metadata":[{
			"ratingKey":"1234","type":"show","title":"Example",
			"Guid":[{"id":"imdb://tt10308554"},{"id":"tmdb://98624"},{"id":"tvdb://376944"}]
		}]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := plex.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.LibrarySectionGetAllWithResponse(t.Context(), "2", &plex.LibrarySectionGetAllParams{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.JSON200 == nil || len(resp.JSON200.MediaContainer.Metadata) != 1 {
		t.Fatalf("unexpected response: %s", resp.Body)
	}

	var got []string
	for _, g := range resp.JSON200.MediaContainer.Metadata[0].GUID {
		got = append(got, g.ID)
	}
	want := []string{"imdb://tt10308554", "tmdb://98624", "tvdb://376944"}
	if !slices.Equal(got, want) {
		t.Errorf("GUID IDs = %q, want %q", got, want)
	}
}
