// Package server serves the generator's page, its JSON API and sprite images.
package server

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/romaine-life/logline"
	"github.com/romaine-life/pokemon-generator/gen"
)

// Config is what the server needs from its environment.
type Config struct {
	Static    fs.FS           // the page files: index.html, app.js, style.css, favicon.svg
	PublicURL string          // the site's origin, for share-card links
	Version   string          // the release this build is
	Log       *logline.Logger // the service's log: stdout, which the cluster keeps
}

type server struct {
	Config
	index *template.Template
}

// New returns the site's handler.
func New(c Config) (http.Handler, error) {
	index, err := template.ParseFS(c.Static, "index.html")
	if err != nil {
		return nil, fmt.Errorf("parse index.html: %w", err)
	}
	s := &server{Config: c, index: index}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("GET /api/monster", s.monster)
	mux.HandleFunc("GET /api/types", s.types)
	mux.HandleFunc("GET /sprite.png", s.sprite)
	mux.HandleFunc("GET /{$}", s.page)
	mux.Handle("GET /static/", http.StripPrefix("/static/", staticFiles(c.Static)))
	return s.logged(mux), nil
}

func (s *server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": s.Version, "generator": gen.Version})
}

// options reads a creature's address from the query. A missing seed gets a fresh random
// one, which the response names so the page can link to it.
func options(q url.Values) (gen.Options, bool, error) {
	o := gen.Options{Seed: gen.NormalizeSeed(q.Get("seed")), Type: q.Get("type")}
	if !gen.ValidType(o.Type) {
		return o, false, fmt.Errorf("unknown type %q", o.Type)
	}
	if st := q.Get("stage"); st != "" {
		n, err := strconv.Atoi(st)
		if err != nil || n < 1 || n > 3 {
			return o, false, fmt.Errorf("stage must be 1, 2 or 3")
		}
		o.Stage = n
	}
	random := o.Seed == ""
	if random {
		o.Seed = randomSeed()
	}
	return o, random, nil
}

func (s *server) monster(w http.ResponseWriter, r *http.Request) {
	o, random, err := options(r.URL.Query())
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.cacheable(w, r, !random)
	writeJSON(w, http.StatusOK, gen.Generate(o))
}

func (s *server) types(w http.ResponseWriter, r *http.Request) {
	s.cacheable(w, r, true)
	writeJSON(w, http.StatusOK, gen.Types())
}

func (s *server) sprite(w http.ResponseWriter, r *http.Request) {
	o, random, err := options(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	scale := 8
	if v := r.URL.Query().Get("scale"); v != "" {
		if scale, err = strconv.Atoi(v); err != nil || scale < 1 || scale > 16 {
			http.Error(w, "scale must be 1 to 16", http.StatusBadRequest)
			return
		}
	}
	var b bytes.Buffer
	m := gen.Generate(o)
	if err := m.Sprite.WritePNG(&b, scale); err != nil {
		http.Error(w, "couldn't draw the sprite", http.StatusInternalServerError)
		return
	}
	s.cacheable(w, r, !random)
	w.Header().Set("Content-Type", "image/png")
	if r.URL.Query().Get("download") != "" {
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", m.Name+".png"))
	}
	w.Write(b.Bytes())
}

// page serves the app. A link to one creature carries its name and sprite in the share-card
// tags, so a pasted link previews the creature.
func (s *server) page(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title, Description, Image, URL, Version string
	}{
		Title:       "Pokémon Generator",
		Description: "Invent a brand-new Pokémon: its sprite, types, stats, moves and evolutions. Type any word to find its creature.",
		Image:       s.PublicURL + "/sprite.png?seed=pokemon-generator&scale=8",
		URL:         s.PublicURL + "/",
		Version:     s.Version,
	}
	q := r.URL.Query()
	if q.Get("seed") != "" {
		if o, _, err := options(q); err == nil {
			m := gen.Generate(o)
			link := url.Values{"seed": {m.Seed}}
			if o.Type != "" {
				link.Set("type", o.Type)
			}
			if o.Stage != 0 {
				link.Set("stage", strconv.Itoa(o.Stage))
			}
			data.Title = fmt.Sprintf("%s · #%04d · Pokémon Generator", m.Name, m.Number)
			data.Description = fmt.Sprintf("The %s Pokémon. %s", m.Category, m.Entry)
			data.Image = s.PublicURL + "/sprite.png?scale=8&v=" + url.QueryEscape(s.Version) + "&" + link.Encode()
			data.URL = s.PublicURL + "/?" + link.Encode()
		}
	}
	var b bytes.Buffer
	if err := s.index.Execute(&b, data); err != nil {
		s.Log.Error("render index", "err", err)
		http.Error(w, "couldn't render the page", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(b.Bytes())
}

// staticFiles serves the page's files, revalidated on each load so a release shows at once.
func staticFiles(files fs.FS) http.Handler {
	h := http.FileServerFS(files)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		h.ServeHTTP(w, r)
	})
}

// cacheable marks a response that its address fully determines: the same seed always gives the
// same creature within a release. The page puts the release in each address it asks for
// (v=), so a new release is a new address and those answers can be kept for a day; an
// address without it may change at the next release, so it is kept for an hour.
func (s *server) cacheable(w http.ResponseWriter, r *http.Request, ok bool) {
	switch {
	case !ok:
		w.Header().Set("Cache-Control", "no-store")
	case r.URL.Query().Get("v") == s.Version:
		w.Header().Set("Cache-Control", "public, max-age=86400")
	default:
		w.Header().Set("Cache-Control", "public, max-age=3600")
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// randomSeed makes a short seed that is easy to read aloud: consonant-vowel pairs and a
// number, like "tavoki42".
func randomSeed() string {
	const consonants, vowels = "bdfgjklmnprstvz", "aeiou"
	b := make([]byte, 8)
	rand.Read(b)
	out := make([]byte, 0, 8)
	for i := 0; i < 3; i++ {
		out = append(out, consonants[int(b[i*2])%len(consonants)], vowels[int(b[i*2+1])%len(vowels)])
	}
	return string(out) + strconv.Itoa(int(b[6])%90+10)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

// logged writes one line per request, except health checks.
func (s *server) logged(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		if r.URL.Path == "/healthz" {
			return
		}
		s.Log.Info("request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery,
			"status", sw.status, "ms", time.Since(start).Milliseconds())
	})
}
