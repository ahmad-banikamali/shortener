package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

func TestResolve_Found(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{
		"abc": {Code: "abc", URL: "https://example.com"},
	}}
	link, err := service.NewResolver(store).Resolve(context.Background(), "abc")
	if err != nil || link.URL != "https://example.com" {
		t.Fatalf("got %v, %v", link, err)
	}
}

func TestResolve_NotFound(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	_, err := service.NewResolver(store).Resolve(context.Background(), "nope")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
