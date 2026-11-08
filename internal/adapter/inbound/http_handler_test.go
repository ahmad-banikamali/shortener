package inbound_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ahmad-banikamali/shortener/internal/adapter/inbound"
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
	h := inbound.NewHTTPHandler(fakeShorten{err: domain.ErrInvalidURL}, fakeResolve{}, fakeStats{err: domain.ErrNotFound}, "http://x")
	mux := http.NewServeMux()
	h.Register(mux)

	req := httptest.NewRequest("POST", "/links", strings.NewReader(`{"url":"bad"}`))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d", rec.Code)
	}
}

type fakeStats struct {
	link domain.Link
	err  error
}

func (f fakeStats) Stats(context.Context, string) (domain.Link, error) { return f.link, f.err }

func TestStats_NotFound_Returns404(t *testing.T) {
	h := inbound.NewHTTPHandler(fakeShorten{}, fakeResolve{}, fakeStats{err: domain.ErrNotFound}, "http://x")
	mux := http.NewServeMux()
	h.Register(mux)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/links/nope/stats", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d", rec.Code)
	}
}
