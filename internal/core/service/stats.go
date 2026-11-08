package service

import (
	"context"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

var _ port.StatsUseCase = (*StatsService)(nil)

type StatsService struct {
	store port.LinkStore
}

func NewStatsService(store port.LinkStore) *StatsService {
	return &StatsService{store: store}
}

// Stats reads a link without counting a visit. Looking at stats must not change them.
func (s *StatsService) Stats(ctx context.Context, code string) (domain.Link, error) {
	if code == "" {
		return domain.Link{}, domain.ErrNotFound
	}
	return s.store.FindByCode(ctx, code)
}
