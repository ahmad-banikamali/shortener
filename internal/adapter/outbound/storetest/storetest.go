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
}
