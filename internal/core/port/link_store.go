package port

import (
	"context"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
)

type LinkStore interface {
	// Save must return domain.ErrCodeTaken if the code already exists.
	Save(ctx context.Context, link domain.Link) error
	// FindByCode must return domain.ErrNotFound if missing.
	FindByCode(ctx context.Context, code string) (domain.Link, error)
}
