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
)

type Link struct {
	Code      string
	URL       string
	CreatedAt time.Time
	Visits    int64
}

// ValidateURL is a business rule, so it lives in the core.
func ValidateURL(raw string) error {
	u, err := url.ParseRequestURI(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return ErrInvalidURL
	}
	return nil
}
