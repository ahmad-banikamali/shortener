package inbound

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
)

type fakeShorten struct{ err error }

func (f fakeShorten) Shorten(_ context.Context, u string) (domain.Link, error) {
	return domain.Link{Code: "abc", URL: u}, f.err
}

type fakeResolve struct{}

func (fakeResolve) Resolve(context.Context, string) (domain.Link, error) {
	return domain.Link{}, domain.ErrNotFound
}

func TestShorten_InvalidURL_Returns400(t *testing.T) {
	h := NewHTTPHandler(fakeShorten{err: domain.ErrInvalidURL}, fakeResolve{}, "http://x")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/links", strings.NewReader(`{"url":"bad"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d", rec.Code)
	}
}
