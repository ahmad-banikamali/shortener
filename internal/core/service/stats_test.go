package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

func TestStats_DoesNotCountVisit(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{
		"abc": {Code: "abc", URL: "https://example.com", Visits: 5},
	}}
	svc := service.NewStatsService(store)

	link, err := svc.Stats(context.Background(), "abc")
	if err != nil || link.Visits != 5 {
		t.Fatalf("got %v, %v", link, err)
	}
	if store.links["abc"].Visits != 5 {
		t.Fatal("stats must not change the visit count")
	}
}

func TestStats_NotFound(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	_, err := service.NewStatsService(store).Stats(context.Background(), "nope")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
