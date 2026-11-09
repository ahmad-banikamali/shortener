package inbound_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/adapter/inbound"
	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

type fakeShorten struct{ err error }

func (f fakeShorten) Shorten(_ context.Context, u string, _ time.Duration) (domain.Link, error) {
	return domain.Link{Code: "abc", URL: u}, f.err
}

type fakeResolve struct {
	link domain.Link
	err  error
}

func (f fakeResolve) Resolve(context.Context, string) (domain.Link, error) { return f.link, f.err }

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

func newMux(s port.ShortenUseCase, r port.ResolveUseCase) *http.ServeMux {
	h := inbound.NewHTTPHandler(s, r, fakeStats{}, "http://x")
	mux := http.NewServeMux()
	h.Register(mux)
	return mux
}

func TestShorten_Success_Returns201(t *testing.T) {
	mux := newMux(fakeShorten{}, fakeResolve{})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("POST", "/links",
		strings.NewReader(`{"url":"https://example.com"}`)))

	if rec.Code != http.StatusCreated || !strings.Contains(rec.Body.String(), `"short":"http://x/abc"`) {
		t.Fatalf("code=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestShorten_InvalidTTL_Returns400(t *testing.T) {
	cases := map[string]string{
		"negative": `{"url":"https://example.com","ttl_seconds":-1}`,
		"overflow": `{"url":"https://example.com","ttl_seconds":18446744074}`,
	}
	for name, body := range cases {
		t.Run(name, func(t *testing.T) {
			mux := newMux(fakeShorten{}, fakeResolve{}) // fake returns no error on purpose

			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, httptest.NewRequest("POST", "/links", strings.NewReader(body)))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("got %d", rec.Code)
			}
		})
	}
}

func TestResolve_Found_Redirects(t *testing.T) {
	mux := newMux(fakeShorten{}, fakeResolve{link: domain.Link{URL: "https://go.dev"}})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/abc", nil))

	if rec.Code != http.StatusFound || rec.Header().Get("Location") != "https://go.dev" {
		t.Fatalf("code=%d location=%q", rec.Code, rec.Header().Get("Location"))
	}
}

func TestResolve_NotFound_Returns404(t *testing.T) {
	mux := newMux(fakeShorten{}, fakeResolve{err: domain.ErrNotFound})

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest("GET", "/nope", nil))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d", rec.Code)
	}
}
