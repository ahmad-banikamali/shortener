package port

import (
	"context"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
)

type ShortenUseCase interface {
	Shorten(ctx context.Context, rawURL string, ttl time.Duration) (domain.Link, error)
}

type ResolveUseCase interface {
	Resolve(ctx context.Context, code string) (domain.Link, error)
}

type StatsUseCase interface {
	Stats(ctx context.Context, code string) (domain.Link, error)
}
