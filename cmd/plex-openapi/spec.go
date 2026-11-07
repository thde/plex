package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

// maxPageSize limits how much of the documentation page is read. The page is
// about 5 MB, most of it being the embedded specification.
const maxPageSize = 64 << 20

// redocState matches the start of the Redoc hydration state assignment,
// e.g. `const __redoc_state = {...};`.
var redocState = regexp.MustCompile(`__redoc_state\s*=\s*`)

// fetchSpec downloads the documentation page at url and returns the embedded
// OpenAPI specification.
func fetchSpec(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "plex-openapi (+https://github.com/thde/plex)")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET %s: unexpected status %s", url, resp.Status)
	}

	page, err := io.ReadAll(io.LimitReader(resp.Body, maxPageSize+1))
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	if len(page) > maxPageSize {
		return nil, fmt.Errorf("GET %s: page exceeds %d bytes", url, maxPageSize)
	}

	spec, err := extractSpec(page)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", url, err)
	}
	return spec, nil
}

// extractSpec returns the OpenAPI specification embedded in a Redoc page as
// indented JSON. The key order of the original document is preserved.
func extractSpec(page []byte) ([]byte, error) {
	loc := redocState.FindIndex(page)
	if loc == nil {
		return nil, errors.New("no __redoc_state found in page")
	}

	// The decoder stops after the first complete JSON value, ignoring the
	// rest of the script.
	var state struct {
		Spec struct {
			Data json.RawMessage `json:"data"`
		} `json:"spec"`
	}
	if err := json.NewDecoder(bytes.NewReader(page[loc[1]:])).Decode(&state); err != nil {
		return nil, fmt.Errorf("decoding __redoc_state: %w", err)
	}

	data := state.Spec.Data
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil, errors.New("__redoc_state contains no spec.data")
	}

	var header struct {
		OpenAPI string `json:"openapi"`
	}
	if err := json.Unmarshal(data, &header); err != nil {
		return nil, fmt.Errorf("spec.data is not an object: %w", err)
	}
	if header.OpenAPI == "" {
		return nil, errors.New("spec.data is not an OpenAPI document: missing openapi version")
	}

	var out bytes.Buffer
	if err := json.Indent(&out, data, "", "  "); err != nil {
		return nil, err
	}
	out.WriteByte('\n')
	return out.Bytes(), nil
}
