package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

type fakeClock struct{ now time.Time }

func (f *fakeClock) Now() time.Time { return f.now }

func TestResolve_Found(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{
		"abc": {Code: "abc", URL: "https://example.com"},
	}}
	clock := &fakeClock{now: time.Now()}
	link, err := service.NewResolver(store, clock).Resolve(context.Background(), "abc")
	if err != nil || link.URL != "https://example.com" {
		t.Fatalf("got %v, %v", link, err)
	}
}

func TestResolve_NotFound(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	clock := &fakeClock{now: time.Now()}
	_, err := service.NewResolver(store, clock).Resolve(context.Background(), "nope")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestResolve_IncrementsVisits(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{
		"abc": {Code: "abc", URL: "https://example.com"},
	}}
	clock := &fakeClock{now: time.Now()}
	svc := service.NewResolver(store, clock)

	svc.Resolve(context.Background(), "abc")
	svc.Resolve(context.Background(), "abc")

	if store.links["abc"].Visits != 2 {
		t.Fatalf("expected 2 visits, got %d", store.links["abc"].Visits)
	}
}

func TestResolve_ExpiredLink_ReturnsNotFound(t *testing.T) {
	start := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := &fakeClock{now: start}
	store := &fakeStore{links: map[string]domain.Link{}}

	shortener := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, clock)
	resolver := service.NewResolver(store, clock)

	shortener.Shorten(context.Background(), "https://example.com", time.Hour)

	clock.now = start.Add(30 * time.Minute) // before expiry
	if _, err := resolver.Resolve(context.Background(), "abc"); err != nil {
		t.Fatalf("expected link to resolve, got %v", err)
	}

	clock.now = start.Add(time.Hour) // exactly at expiry
	if _, err := resolver.Resolve(context.Background(), "abc"); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after expiry, got %v", err)
	}
}
