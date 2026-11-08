package port

import (
	"context"

	"github.com/ahmad-banikamali/internal/core/domain"
)

type ShortenUseCase interface {
	Shorten(ctx context.Context, rawURL string) (domain.Link, error)
}

type ResolveUseCase interface {
	Resolve(ctx context.Context, code string) (domain.Link, error)
}
