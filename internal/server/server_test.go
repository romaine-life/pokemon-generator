package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/romaine-life/logline"
	"github.com/romaine-life/pokemon-generator/web"
)

func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()
	h, err := New(Config{Static: web.Static(), PublicURL: "https://pg.test", Version: "test",
		Log: logline.New(io.Discard)})
	if err != nil {
		t.Fatal(err)
	}
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	return srv
}

func get(t *testing.T, url string) (*http.Response, string) {
	t.Helper()
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return res, string(b)
}

func TestMonsterBySeedIsCacheable(t *testing.T) {
	srv := newTestServer(t)
	res, body := get(t, srv.URL+"/api/monster?seed=Pikachu&type=water&stage=1&v=test")
	if res.StatusCode != 200 || res.Header.Get("Cache-Control") != "public, max-age=86400" {
		t.Fatalf("status %d, cache %q", res.StatusCode, res.Header.Get("Cache-Control"))
	}
	// Without this release's version the answer may change at the next release.
	for _, v := range []string{"", "&v=older"} {
		if res, _ := get(t, srv.URL+"/api/monster?seed=Pikachu"+v); res.Header.Get("Cache-Control") != "public, max-age=3600" {
			t.Fatalf("v=%q: cache %q", v, res.Header.Get("Cache-Control"))
		}
	}
	var m struct {
		Seed, TypeLock string
		Stage          int
	}
	if err := json.Unmarshal([]byte(body), &m); err != nil {
		t.Fatal(err)
	}
	if m.Seed != "pikachu" || m.TypeLock != "water" || m.Stage != 1 {
		t.Fatalf("got %+v", m)
	}
}

func TestRandomMonsterNamesItsSeed(t *testing.T) {
	srv := newTestServer(t)
	res, body := get(t, srv.URL+"/api/monster")
	if res.StatusCode != 200 || res.Header.Get("Cache-Control") != "no-store" {
		t.Fatalf("status %d, cache %q", res.StatusCode, res.Header.Get("Cache-Control"))
	}
	var m struct{ Seed string }
	if err := json.Unmarshal([]byte(body), &m); err != nil || len(m.Seed) != 8 {
		t.Fatalf("random seed %q (%v)", m.Seed, err)
	}
}

func TestBadAddressesAreRefused(t *testing.T) {
	srv := newTestServer(t)
	for _, path := range []string{
		"/api/monster?type=plasma",
		"/api/monster?seed=a&stage=4",
		"/api/monster?seed=a&stage=two",
		"/sprite.png?seed=a&scale=99",
	} {
		if res, body := get(t, srv.URL+path); res.StatusCode != http.StatusBadRequest || body == "" {
			t.Errorf("%s: status %d, body %q", path, res.StatusCode, body)
		}
	}
	if res, _ := get(t, srv.URL+"/nope"); res.StatusCode != http.StatusNotFound {
		t.Errorf("/nope: status %d", res.StatusCode)
	}
}

func TestSharedLinkPreviewsTheCreature(t *testing.T) {
	srv := newTestServer(t)
	_, body := get(t, srv.URL+"/?seed=Nelson&stage=1")
	for _, want := range []string{
		`<meta property="og:image" content="https://pg.test/sprite.png?scale=8&amp;v=test&amp;seed=nelson&amp;stage=1">`,
		`<meta property="og:url" content="https://pg.test/?seed=nelson&amp;stage=1">`,
		"· Pokémon Generator</title>",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page lacks %s", want)
		}
	}
}

func TestSpriteDownload(t *testing.T) {
	srv := newTestServer(t)
	res, body := get(t, srv.URL+"/sprite.png?seed=nelson&scale=2&download=1")
	if res.Header.Get("Content-Type") != "image/png" || !strings.HasPrefix(body, "\x89PNG") {
		t.Fatalf("content type %q", res.Header.Get("Content-Type"))
	}
	if cd := res.Header.Get("Content-Disposition"); !strings.HasPrefix(cd, "attachment; filename=") || !strings.HasSuffix(cd, `.png"`) {
		t.Fatalf("disposition %q", cd)
	}
}

// Each request but a health check is one standard line on the log.
func TestRequestsAreLogged(t *testing.T) {
	var logs bytes.Buffer
	h, err := New(Config{Static: web.Static(), PublicURL: "https://pg.test", Version: "test", Log: logline.New(&logs)})
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"/healthz", "/api/monster?seed=nelson&type=fire", "/api/types"} {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", path, nil))
	}
	want := []map[string]string{
		{"method": "GET", "path": "/api/monster", "query": "seed=nelson&type=fire", "status": "200"},
		{"method": "GET", "path": "/api/types", "query": "", "status": "200"},
	}
	lines := strings.Split(strings.TrimSuffix(logs.String(), "\n"), "\n")
	if len(lines) != len(want) {
		t.Fatalf("want %d lines, health checks unlogged: %q", len(want), logs.String())
	}
	for i, line := range lines {
		e, err := logline.Parse(line)
		if err != nil {
			t.Fatalf("not the standard line (%v): %q", err, line)
		}
		got := map[string]string{}
		for _, f := range e.Fields {
			got[f.Key] = f.Value
		}
		if e.Level != logline.Info || e.Message != "request" || got["ms"] == "" {
			t.Errorf("line %d: %q", i, line)
		}
		for k, v := range want[i] {
			if got[k] != v {
				t.Errorf("line %d: %s = %q, want %q", i, k, got[k], v)
			}
		}
	}
}
