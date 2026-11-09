package domain_test

import (
	"testing"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
)

func TestLink_IsExpired(t *testing.T) {
	exp := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

	cases := []struct {
		name      string
		expiresAt *time.Time
		now       time.Time
		want      bool
	}{
		{"no expiry never expires", nil, exp.Add(100 * 365 * 24 * time.Hour), false},
		{"1ns before expiry", &exp, exp.Add(-time.Nanosecond), false},
		{"exactly at expiry", &exp, exp, true},
		{"after expiry", &exp, exp.Add(time.Second), true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := domain.Link{ExpiresAt: c.expiresAt}
			if got := l.IsExpired(c.now); got != c.want {
				t.Fatalf("IsExpired = %v, want %v", got, c.want)
			}
		})
	}
}
