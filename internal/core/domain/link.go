package domain

import (
	"errors"
	"net/url"
	"time"
)

var (
	ErrInvalidURL = errors.New("invalid url")
	ErrNotFound   = errors.New("link not found")
	ErrCodeTaken  = errors.New("code already taken")
	ErrInvalidTTL = errors.New("invalid ttl")
)

const MaxTTL = 10 * 365 * 24 * time.Hour

type Link struct {
	Code      string
	URL       string
	CreatedAt time.Time
	Visits    int64
	ExpiresAt *time.Time // nil = never expires
}

// ValidateURL is a business rule, so it lives in the core.
func ValidateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}

func (l Link) IsExpired(now time.Time) bool {
	return l.ExpiresAt != nil && !now.Before(*l.ExpiresAt)
}
