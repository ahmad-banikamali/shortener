package main

import (
	"context"
	"log"
	"os"

	"github.com/ahmad-banikamali/shortener/internal/adapter/inbound"
	"github.com/ahmad-banikamali/shortener/internal/adapter/outbound"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

func main() {
	store, err := outbound.NewSQLiteStore("shortener.db")
	if err != nil {
		log.Fatal(err)
	}

	codes := outbound.RandomCodeGenerator{Length: 6}
	clock := outbound.SystemClock{}
	cli := inbound.NewCLI(
		service.NewShortener(store, codes, clock),
		service.NewResolver(store, clock),
		service.NewStatsService(store),
		"http://localhost:8080",
		os.Stdout, os.Stderr,
	)

	exitCode := cli.Run(context.Background(), os.Args[1:])
	store.Close() // os.Exit skips defers, so close explicitly
	os.Exit(exitCode)
}
