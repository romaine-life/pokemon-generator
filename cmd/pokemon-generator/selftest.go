package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/romaine-life/logline"
	"github.com/romaine-life/pokemon-generator/internal/server"
	"github.com/romaine-life/pokemon-generator/web"
)

// runSelfTest serves the built image's own files and checks each route answers as the page
// expects.
func runSelfTest() error {
	h, err := server.New(server.Config{Static: web.Static(), PublicURL: "https://example.test", Version: version,
		Log: logline.New(io.Discard)})
	if err != nil {
		return err
	}
	srv := httptest.NewServer(h)
	defer srv.Close()

	checks := []struct{ path, ctype, contains string }{
		{"/healthz", "application/json", `"status":"ok"`},
		{"/", "text/html", "/static/app.js"},
		{"/?seed=self-test", "text/html", "/sprite.png?scale=8&amp;v="},
		{"/static/app.js", "text/javascript", "drawSprite"},
		{"/static/style.css", "text/css", ".card"},
		{"/static/favicon.svg", "image/svg+xml", "<svg"},
		{"/api/types", "application/json", `"id":"dragon"`},
		{"/sprite.png?seed=self-test", "image/png", "\x89PNG"},
	}
	for _, c := range checks {
		res, err := http.Get(srv.URL + c.path)
		if err != nil {
			return fmt.Errorf("GET %s: %w", c.path, err)
		}
		body, _ := io.ReadAll(res.Body)
		res.Body.Close()
		if res.StatusCode != http.StatusOK {
			return fmt.Errorf("GET %s: status %d", c.path, res.StatusCode)
		}
		if ct := res.Header.Get("Content-Type"); !strings.HasPrefix(ct, c.ctype) {
			return fmt.Errorf("GET %s: content type %q, want %s", c.path, ct, c.ctype)
		}
		if !strings.Contains(string(body), c.contains) {
			return fmt.Errorf("GET %s: body lacks %q", c.path, c.contains)
		}
	}

	res, err := http.Get(srv.URL + "/api/monster?seed=self-test&type=fire")
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var m struct {
		Name   string
		Types  []struct{ ID string }
		Sprite struct{ Pixels []int }
	}
	if err := json.NewDecoder(res.Body).Decode(&m); err != nil {
		return fmt.Errorf("decode monster: %w", err)
	}
	if m.Name == "" || len(m.Types) == 0 || m.Types[0].ID != "fire" || len(m.Sprite.Pixels) != 40*40 {
		return fmt.Errorf("monster for seed self-test is incomplete: %+v", m.Name)
	}
	return nil
}
