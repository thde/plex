package plex_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"thde.io/plex"
)

// newServer returns a test server that records the headers of the last
// request and responds with a minimal identity document.
func newServer(t *testing.T) (*httptest.Server, *http.Header) {
	t.Helper()

	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"MediaContainer":{"size":0,"claimed":true,"machineIdentifier":"abc123","version":"1.42.0"}}`))
	}))
	t.Cleanup(srv.Close)

	return srv, &got
}

func TestNew(t *testing.T) {
	srv, got := newServer(t)

	client, err := plex.New(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.GetIdentityWithResponse(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if want := "application/json"; got.Get("Accept") != want {
		t.Errorf("Accept = %q, want %q", got.Get("Accept"), want)
	}
	if resp.JSON200 == nil {
		t.Fatalf("JSON200 = nil, body: %s", resp.Body)
	}
	if mc := resp.JSON200.MediaContainer; mc.MachineIdentifier != "abc123" || !mc.Claimed {
		t.Errorf("unexpected MediaContainer: %+v", mc)
	}
}

func TestOptions(t *testing.T) {
	tests := []struct {
		name   string
		option plex.ClientOption
		want   map[string]string
	}{
		{
			name:   "token",
			option: plex.WithToken("my-token"),
			want:   map[string]string{"X-Plex-Token": "my-token"},
		},
		{
			name:   "client identifier",
			option: plex.WithClientIdentifier("my-client"),
			want:   map[string]string{"X-Plex-Client-Identifier": "my-client"},
		},
		{
			name:   "product",
			option: plex.WithProduct("my-app", "1.0"),
			want:   map[string]string{"X-Plex-Product": "my-app", "X-Plex-Version": "1.0"},
		},
		{
			name:   "user agent",
			option: plex.WithUserAgent("my-agent/1.0"),
			want:   map[string]string{"User-Agent": "my-agent/1.0"},
		},
		{
			name: "accept override",
			option: plex.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
				req.Header.Set("Accept", "application/xml")
				return nil
			}),
			want: map[string]string{"Accept": "application/xml"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, got := newServer(t)

			client, err := plex.New(srv.URL, tt.option)
			if err != nil {
				t.Fatal(err)
			}

			if _, err := client.GetIdentityWithResponse(t.Context()); err != nil {
				t.Fatal(err)
			}

			for key, want := range tt.want {
				if v := got.Get(key); v != want {
					t.Errorf("%s = %q, want %q", key, v, want)
				}
			}
		})
	}
}

func TestWithClientIdentifier_ParameterTakesPrecedence(t *testing.T) {
	srv, got := newServer(t)

	client, err := plex.New(srv.URL, plex.WithClientIdentifier("default-client"))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.TranscodeSubtitlesWithResponse(t.Context(), plex.TranscodeSubtitlesParamsTranscodeTypeSubtitles, &plex.TranscodeSubtitlesParams{
		XPlexClientIdentifier: "param-client",
	})
	if err != nil {
		t.Fatal(err)
	}

	if want := "param-client"; got.Get("X-Plex-Client-Identifier") != want {
		t.Errorf("X-Plex-Client-Identifier = %q, want %q", got.Get("X-Plex-Client-Identifier"), want)
	}
}
