package plex_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"thde.io/plex"
)

// TestHistoryParentFields decodes a play of an episode. Plex Media Server
// links it to its season and show with parentKey and grandparentKey, which
// the specification leaves out.
func TestHistoryParentFields(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"MediaContainer":{"size":1,"totalSize":1,"offset":0,"Metadata":[{
			"historyKey":"/status/sessions/history/508","key":"/library/metadata/300","ratingKey":"300",
			"librarySectionID":"2","parentKey":"/library/metadata/200","grandparentKey":"/library/metadata/100",
			"title":"Episode","grandparentTitle":"Show","type":"episode","index":10,"parentIndex":1,
			"viewedAt":1793966686,"accountID":1,"deviceID":417
		}]}}`))
	}))
	t.Cleanup(srv.Close)

	client, err := plex.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.StatusGetHistoryAllWithResponse(t.Context(), &plex.StatusGetHistoryAllParams{})
	if err != nil {
		t.Fatal(err)
	}
	if resp.JSON200 == nil || len(resp.JSON200.MediaContainer.Metadata) != 1 {
		t.Fatalf("unexpected response: %s", resp.Body)
	}

	got := resp.JSON200.MediaContainer.Metadata[0]
	if got.ParentKey != "/library/metadata/200" {
		t.Errorf("ParentKey = %q", got.ParentKey)
	}
	if got.GrandparentKey != "/library/metadata/100" {
		t.Errorf("GrandparentKey = %q", got.GrandparentKey)
	}
	if got.GrandparentTitle != "Show" {
		t.Errorf("GrandparentTitle = %q", got.GrandparentTitle)
	}
	if got.Index != 10 || got.ParentIndex != 1 {
		t.Errorf("Index, ParentIndex = %d, %d, want 10, 1", got.Index, got.ParentIndex)
	}
}
