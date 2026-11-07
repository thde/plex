// Package plex provides a client for the Plex Media Server API.
package plex

//go:generate go run ./cmd/plex-openapi -o openapi.json
//go:generate go tool oapi-codegen -config oapi-codegen.yml openapi.json

import (
	"context"
	"net/http"
)

const (
	headerAccept           = "Accept"
	headerUserAgent        = "User-Agent"
	headerToken            = "X-Plex-Token"
	headerClientIdentifier = "X-Plex-Client-Identifier"
	headerProduct          = "X-Plex-Product"
	headerVersion          = "X-Plex-Version"

	mimeJSON = "application/json"
)

// WithToken returns a [ClientOption] that authenticates requests with the
// provided Plex token using the X-Plex-Token header.
func WithToken(token string) ClientOption {
	return withHeader(headerToken, token)
}

// WithClientIdentifier returns a [ClientOption] that sets the
// X-Plex-Client-Identifier header. The identifier should be unique and stable
// per client installation.
//
// Operations that take the client identifier as a parameter use the value
// from the parameter instead, if set.
func WithClientIdentifier(id string) ClientOption {
	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if req.Header.Get(headerClientIdentifier) == "" {
			req.Header.Set(headerClientIdentifier, id)
		}
		return nil
	})
}

// WithProduct returns a [ClientOption] that sets the X-Plex-Product and
// X-Plex-Version headers identifying the client application.
func WithProduct(product, version string) ClientOption {
	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set(headerProduct, product)
		req.Header.Set(headerVersion, version)
		return nil
	})
}

// WithUserAgent returns a [ClientOption] that sets the User-Agent header.
func WithUserAgent(agent string) ClientOption {
	return withHeader(headerUserAgent, agent)
}

func withHeader(key, value string) ClientOption {
	return WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set(key, value)
		return nil
	})
}

// New creates a new Plex Media Server client for the server at baseURL, e.g.
// "http://192.168.1.10:32400" or "https://1-2-3-4.<id>.plex.direct:32400".
//
// Plex Media Server responds with XML unless JSON is requested, so the client
// sends "Accept: application/json" by default. The generated response parsers
// only handle JSON.
func New(baseURL string, opts ...ClientOption) (*ClientWithResponses, error) {
	opts = append([]ClientOption{withHeader(headerAccept, mimeJSON)}, opts...)
	return NewClientWithResponses(baseURL, opts...)
}
