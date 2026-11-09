package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/service"
)

type fakeStore struct{ links map[string]domain.Link }

func (f *fakeStore) IncrementVisits(_ context.Context, code string) error {
	l, ok := f.links[code]
	if !ok {
		return domain.ErrNotFound
	}
	l.Visits++
	f.links[code] = l
	return nil
}

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
	clock := &fakeClock{}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, clock)

	link, err := svc.Shorten(context.Background(), "https://example.com", time.Duration(5)*time.Second)
	if err != nil || link.Code != "abc" {
		t.Fatalf("got %v, %v", link, err)
	}
}

func TestShorten_InvalidURL(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	clock := &fakeClock{now: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, clock)

	_, err := svc.Shorten(context.Background(), "not a url", time.Duration(5)*time.Second)
	if !errors.Is(err, domain.ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got %v", err)
	}
}

func TestShorten_RetriesOnCollision(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{"aaa": {Code: "aaa"}}}
	clock := &fakeClock{now: time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"aaa", "bbb"}}, clock)

	link, err := svc.Shorten(context.Background(), "https://example.com", time.Duration(5)*time.Second)
	if err != nil || link.Code != "bbb" {
		t.Fatalf("expected code bbb, got %v, %v", link, err)
	}
}

func TestShorten_NegativeTTL(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, &fakeClock{now: time.Now()})

	_, err := svc.Shorten(context.Background(), "https://example.com", -time.Second)
	if !errors.Is(err, domain.ErrInvalidTTL) {
		t.Fatalf("expected ErrInvalidTTL, got %v", err)
	}
}

func TestShorten_SetsExpiry(t *testing.T) {
	now := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{links: map[string]domain.Link{}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, &fakeClock{now: now})

	link, err := svc.Shorten(context.Background(), "https://example.com", time.Hour)

	want := now.Add(time.Hour)
	if err != nil || link.ExpiresAt == nil || !link.ExpiresAt.Equal(want) {
		t.Fatalf("expected expiry %v, got %v (err %v)", want, link.ExpiresAt, err)
	}
}

func TestShorten_ZeroTTL_NeverExpires(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, &fakeClock{now: time.Now()})

	link, err := svc.Shorten(context.Background(), "https://example.com", 0)

	if err != nil || link.ExpiresAt != nil {
		t.Fatalf("expected nil expiry, got %v (err %v)", link.ExpiresAt, err)
	}
}

func TestShorten_TTLAboveMax(t *testing.T) {
	store := &fakeStore{links: map[string]domain.Link{}}
	svc := service.NewShortener(store, &seqCodes{codes: []string{"abc"}}, &fakeClock{now: time.Now()})

	_, err := svc.Shorten(context.Background(), "https://example.com", domain.MaxTTL+time.Second)

	if !errors.Is(err, domain.ErrInvalidTTL) {
		t.Fatalf("expected ErrInvalidTTL, got %v", err)
	}
}
