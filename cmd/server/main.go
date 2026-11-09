package main

import (
	"log"
	"net/http"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/adapter/inbound"
	"github.com/ahmad-banikamali/shortener/internal/adapter/outbound"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

func main() {
	// Outbound adapters
	store, err := outbound.NewSQLiteStore("shortener.db")
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	codes := outbound.RandomCodeGenerator{Length: 6}

	// Core
	clock := outbound.SystemClock{}
	shortener := service.NewShortener(store, codes, clock)
	resolver := service.NewResolver(store, clock)
	stats := service.NewStatsService(store)

	handler := inbound.NewHTTPHandler(shortener, resolver, stats, "http://localhost:8080")
	mux := http.NewServeMux()
	handler.Register(mux)

	log.Println("listening on :8080")
	srv := &http.Server{Addr: ":8080", Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
