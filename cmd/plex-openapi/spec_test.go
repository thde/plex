package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const page = `<!DOCTYPE html><html><body>
<div id="redoc"></div>
<script>
const __redoc_state = {"menu":{"activeItemIdx":-1},"spec":{"data":{"openapi":"3.1.0","info":{"title":"Plex Media Server","version":"1"},"paths":{"/z":{},"/a":{}}}},"searchIndex":{},"options":{}};
var container = document.getElementById('redoc');
Redoc.hydrate(__redoc_state, container);
</script></body></html>`

func TestExtractSpec(t *testing.T) {
	got, err := extractSpec([]byte(page))
	if err != nil {
		t.Fatal(err)
	}

	var spec struct {
		OpenAPI string         `json:"openapi"`
		Paths   map[string]any `json:"paths"`
	}
	if err := json.Unmarshal(got, &spec); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if spec.OpenAPI != "3.1.0" || len(spec.Paths) != 2 {
		t.Errorf("unexpected spec: %+v", spec)
	}

	// Key order of the source document must be preserved.
	if z, a := strings.Index(string(got), `"/z"`), strings.Index(string(got), `"/a"`); z > a {
		t.Errorf("key order not preserved:\n%s", got)
	}
	if !strings.HasSuffix(string(got), "}\n") {
		t.Errorf("output should end with a newline")
	}
}

func TestExtractSpecErrors(t *testing.T) {
	tests := map[string]string{
		"no state":     `<html><script>var x = {};</script></html>`,
		"invalid json": `<script>const __redoc_state = {"spec": </script>`,
		"no data":      `<script>const __redoc_state = {"spec":{}};</script>`,
		"null data":    `<script>const __redoc_state = {"spec":{"data":null}};</script>`,
		"not openapi":  `<script>const __redoc_state = {"spec":{"data":{"swagger":"2.0"}}};</script>`,
		"not object":   `<script>const __redoc_state = {"spec":{"data":[1]}};</script>`,
	}
	for name, page := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := extractSpec([]byte(page)); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestFetchSpec(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/pms/" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(page))
	}))
	defer srv.Close()

	if _, err := fetchSpec(t.Context(), srv.URL+"/pms/"); err != nil {
		t.Fatal(err)
	}
	if _, err := fetchSpec(t.Context(), srv.URL+"/missing"); err == nil {
		t.Error("expected error for non-200 response, got nil")
	}
}
