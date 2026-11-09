package service

import (
	"context"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

var _ port.ResolveUseCase = (*Resolver)(nil)

type Resolver struct {
	store port.LinkStore
	clock port.Clock
}

func NewResolver(store port.LinkStore, clock port.Clock) *Resolver {
	return &Resolver{store: store, clock: clock}
}

func (r *Resolver) Resolve(ctx context.Context, code string) (domain.Link, error) {
	if code == "" {
		return domain.Link{}, domain.ErrNotFound
	}
	link, err := r.store.FindByCode(ctx, code)
	if err != nil {
		return domain.Link{}, err
	}
	if link.IsExpired(r.clock.Now()) {
		return domain.Link{}, domain.ErrNotFound
	}
	// Business decision: a failed counter must not break the redirect.
	_ = r.store.IncrementVisits(ctx, code)
	return link, nil
}
