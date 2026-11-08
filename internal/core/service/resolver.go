package service

import (
	"context"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

var _ port.ResolveUseCase = (*Resolver)(nil)

type Resolver struct {
	store port.LinkStore
}

func NewResolver(store port.LinkStore) *Resolver {
	return &Resolver{store: store}
}

func (r *Resolver) Resolve(ctx context.Context, code string) (domain.Link, error) {
	if code == "" {
		return domain.Link{}, domain.ErrNotFound
	}
	return r.store.FindByCode(ctx, code)
}
