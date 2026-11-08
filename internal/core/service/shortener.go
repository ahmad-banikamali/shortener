package service

import (
	"context"
	"errors"
	"time"

	"github.com/ahmad-banikamali/internal/core/domain"
	"github.com/ahmad-banikamali/internal/core/port"
)

var _ port.ShortenUseCase = (*Shortener)(nil)

type Shortener struct {
	store port.LinkStore
	codes port.CodeGenerator
}

func NewShortener(store port.LinkStore, codes port.CodeGenerator) *Shortener {
	return &Shortener{store: store, codes: codes}
}

func (s *Shortener) Shorten(ctx context.Context, rawURL string) (domain.Link, error) {
	if err := domain.ValidateURL(rawURL); err != nil {
		return domain.Link{}, err
	}

	for i := 0; i < 3; i++ { // retry on code collision
		code, err := s.codes.Generate()
		if err != nil {
			return domain.Link{}, err
		}
		link := domain.Link{Code: code, URL: rawURL, CreatedAt: time.Now()}

		err = s.store.Save(ctx, link)
		if err == nil {
			return link, nil
		}
		if !errors.Is(err, domain.ErrCodeTaken) {
			return domain.Link{}, err
		}
	}
	return domain.Link{}, errors.New("could not generate a unique code")
}
