package main

import (
	"log"
	"net/http"

	"github.com/ahmad-banikamali/internal/adapter/inbound"
	"github.com/ahmad-banikamali/internal/adapter/outbound"
	"github.com/ahmad-banikamali/internal/core/service"
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
	shortener := service.NewShortener(store, codes)
	resolver := service.NewResolver(store)

	// Inbound adapter
	handler := inbound.NewHTTPHandler(shortener, resolver, "http://localhost:8080")
	mux := http.NewServeMux()
	handler.Register(mux)

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
