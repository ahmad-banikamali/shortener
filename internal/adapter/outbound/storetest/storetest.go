package storetest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

func Run(t *testing.T, newStore func(t *testing.T) port.LinkStore) {
	ctx := context.Background()
	link := domain.Link{Code: "abc", URL: "https://example.com", CreatedAt: time.Now()}

	t.Run("save then find", func(t *testing.T) {
		s := newStore(t)
		if err := s.Save(ctx, link); err != nil {
			t.Fatal(err)
		}
		got, err := s.FindByCode(ctx, "abc")
		if err != nil || got.URL != link.URL || got.Code != link.Code {
			t.Fatalf("got %v, %v", got, err)
		}
	})

	t.Run("duplicate code", func(t *testing.T) {
		s := newStore(t)
		_ = s.Save(ctx, link)
		if err := s.Save(ctx, link); !errors.Is(err, domain.ErrCodeTaken) {
			t.Fatalf("expected ErrCodeTaken, got %v", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		s := newStore(t)
		if _, err := s.FindByCode(ctx, "missing"); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("increment visits", func(t *testing.T) {
		s := newStore(t)
		_ = s.Save(ctx, link)
		_ = s.IncrementVisits(ctx, "abc")
		_ = s.IncrementVisits(ctx, "abc")
		got, _ := s.FindByCode(ctx, "abc")
		if got.Visits != 2 {
			t.Fatalf("expected 2 visits, got %d", got.Visits)
		}
	})

	t.Run("increment missing code", func(t *testing.T) {
		s := newStore(t)
		if err := s.IncrementVisits(ctx, "missing"); !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("expiry round trip", func(t *testing.T) {
		s := newStore(t)
		exp := time.Now().Add(time.Hour).Truncate(time.Second)
		withExp := domain.Link{Code: "exp", URL: "https://example.com", CreatedAt: time.Now(), ExpiresAt: &exp}
		if err := s.Save(ctx, withExp); err != nil {
			t.Fatal(err)
		}
		got, err := s.FindByCode(ctx, "exp")
		if err != nil || got.ExpiresAt == nil || !got.ExpiresAt.Equal(exp) {
			t.Fatalf("expiry lost: %v, %v", got.ExpiresAt, err)
		}

		// and a link without expiry must come back with nil
		if err := s.Save(ctx, domain.Link{Code: "never", URL: "https://example.com", CreatedAt: time.Now()}); err != nil {
			t.Fatal(err)
		}
		got, _ = s.FindByCode(ctx, "never")
		if got.ExpiresAt != nil {
			t.Fatalf("expected nil expiry, got %v", got.ExpiresAt)
		}
	})
}
