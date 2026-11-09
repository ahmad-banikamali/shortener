package service

import (
	"context"
	"errors"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

var _ port.ShortenUseCase = (*Shortener)(nil)

type Shortener struct {
	store port.LinkStore
	codes port.CodeGenerator
	clock port.Clock
}

func NewShortener(store port.LinkStore, codes port.CodeGenerator, clock port.Clock) *Shortener {
	return &Shortener{store: store, codes: codes, clock: clock}
}

func (s *Shortener) Shorten(ctx context.Context, rawURL string, ttl time.Duration) (domain.Link, error) {
	if err := domain.ValidateURL(rawURL); err != nil {
		return domain.Link{}, err
	}
	if ttl < 0 || ttl > domain.MaxTTL {
		return domain.Link{}, domain.ErrInvalidTTL
	}

	now := s.clock.Now()
	var expires *time.Time
	if ttl > 0 {
		t := now.Add(ttl)
		expires = &t
	}

	for range 3 {
		code, err := s.codes.Generate()
		if err != nil {
			return domain.Link{}, err
		}
		link := domain.Link{Code: code, URL: rawURL, CreatedAt: now, ExpiresAt: expires}

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
