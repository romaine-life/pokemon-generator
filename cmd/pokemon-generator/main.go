// pokemon-generator serves https://pokemon-generator.romaine.life: a page that invents a new
// creature from any seed, its JSON API and its sprite images.
//
// Environment:
//
//	PORT        listen port (8080)
//	PUBLIC_URL  the site's origin, for share-card links (https://pokemon-generator.romaine.life)
//
// Flags: --self-test checks the built image's routes and exits (the web-app suite's image
// check).
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/romaine-life/logline"
	"github.com/romaine-life/pokemon-generator/internal/server"
	"github.com/romaine-life/pokemon-generator/web"
)

// version is the release, stamped by the image build.
var version = "dev"

func main() {
	selfTest := flag.Bool("self-test", false, "check the built image's routes and exit")
	flag.Parse()
	log := logline.New(os.Stdout)

	if *selfTest {
		if err := runSelfTest(); err != nil {
			fmt.Fprintf(os.Stderr, "pokemon-generator self-test failed: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("pokemon-generator self-test passed")
		return
	}
	if err := run(log); err != nil {
		log.Error("pokemon-generator stopped", "err", err)
		os.Exit(1)
	}
}

func env(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}

func run(log *logline.Logger) error {
	h, err := server.New(server.Config{
		Static:    web.Static(),
		PublicURL: env("PUBLIC_URL", "https://pokemon-generator.romaine.life"),
		Version:   version,
		Log:       log,
	})
	if err != nil {
		return err
	}
	srv := &http.Server{Addr: ":" + env("PORT", "8080"), Handler: h, ReadHeaderTimeout: 10 * time.Second}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errc := make(chan error, 1)
	go func() { errc <- srv.ListenAndServe() }()
	log.Info("listening", "addr", srv.Addr, "version", version)
	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}
	// A rollout's SIGTERM: finish what's in flight, then exit.
	shutdown, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdown); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
