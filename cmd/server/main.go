package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/adapter/inbound"
	"github.com/ahmad-banikamali/shortener/internal/adapter/outbound"
	"github.com/ahmad-banikamali/shortener/internal/config"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg := config.Load()

	store, err := outbound.NewSQLiteStore(cfg.DBPath)
	if err != nil {
		return err
	}
	defer store.Close() // now actually runs

	codes := outbound.RandomCodeGenerator{Length: 6}
	clock := outbound.SystemClock{}

	handler := inbound.NewHTTPHandler(
		service.NewShortener(store, codes, clock),
		service.NewResolver(store, clock),
		service.NewStatsService(store),
		cfg.BaseURL,
	)
	mux := http.NewServeMux()
	handler.Register(mux)

	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}

	// ctx is cancelled when the process receives Ctrl+C or SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		log.Println("listening on", srv.Addr)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		return err // failed to start, e.g. port already in use
	case <-ctx.Done():
		log.Println("shutting down...")
	}

	// Give in-flight requests up to 10 seconds to finish.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}
