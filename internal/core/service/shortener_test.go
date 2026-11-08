package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ahmad-banikamali/internal/core/domain"
	"github.com/ahmad-banikamali/internal/core/service"
)

type fakeStore struct{ links map[string]domain.Link }

func (f *fakeStore) Save(_ context.Context, l domain.Link) error {
	if _, ok := f.links[l.Code]; ok {
		return domain.ErrCodeTaken
	}
	f.links[l.Code] = l
	return nil
}

func (f *fakeStore) FindByCode(_ context.Context, code string) (domain.Link, error) {
	l, ok := f.links[code]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}
	return l, nil
}

type seqCodes struct {
	codes []string
	i     int
}

func (s *seqCodes) Generate() (string, error) {
	c := s.codes[s.i]
	s.i++
	return c, nil
}

func TestShorten_Success(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}})

	link, err := svc.Shorten(context.Background(), "https://example.com")
	if err != nil || link.Code != "abc" {
		t.Fatalf("got %v, %v", link, err)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}})

	_, err := svc.Shorten(context.Background(), "not a url")
	if !errors.Is(err, domain.ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestShorten_RetriesOnCollision(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{"aaa": {Code: "aaa"}}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"aaa", "bbb"}})

	link, err := svc.Shorten(context.Background(), "https://example.com")
	if err != nil || link.Code != "bbb" {
		t.Fatalf("expected code bbb, got %v, %v", link, err)
	}
}
