package inbound

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

type HTTPHandler struct {
	shorten port.ShortenUseCase
	resolve port.ResolveUseCase
	stats   port.StatsUseCase
	baseURL string
}

func NewHTTPHandler(s port.ShortenUseCase, r port.ResolveUseCase, st port.StatsUseCase, baseURL string) *HTTPHandler {
	return &HTTPHandler{shorten: s, resolve: r, stats: st, baseURL: baseURL}
}

func (h *HTTPHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /links", h.handleShorten)
	mux.HandleFunc("GET /links/{code}/stats", h.handleStats)
	mux.HandleFunc("GET /{code}", h.handleResolve)
}

type shortenRequest struct {
	URL        string `json:"url"`
	TTLSeconds int64  `json:"ttl_seconds"`
}

type shortenResponse struct {
	Code  string `json:"code"`
	Short string `json:"short"`
}

func (h *HTTPHandler) handleShorten(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MB limit

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if req.TTLSeconds < 0 || req.TTLSeconds > int64(domain.MaxTTL/time.Second) {
		writeError(w, domain.ErrInvalidTTL)
		return
	}

	link, err := h.shorten.Shorten(r.Context(), req.URL, time.Duration(req.TTLSeconds)*time.Second)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shortenResponse{
		Code:  link.Code,
		Short: h.baseURL + "/" + link.Code,
	})
}

func (h *HTTPHandler) handleResolve(w http.ResponseWriter, r *http.Request) {
	link, err := h.resolve.Resolve(r.Context(), r.PathValue("code"))
	if err != nil {
		writeError(w, err)
		return
	}
	http.Redirect(w, r, link.URL, http.StatusFound)
}

type statsResponse struct {
	Code      string    `json:"code"`
	URL       string    `json:"url"`
	Visits    int64     `json:"visits"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *HTTPHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	link, err := h.stats.Stats(r.Context(), r.PathValue("code"))
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statsResponse{
		Code: link.Code, URL: link.URL, Visits: link.Visits, CreatedAt: link.CreatedAt,
	})
}

// The adapter translates domain errors into HTTP status codes.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidURL):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrInvalidTTL):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
