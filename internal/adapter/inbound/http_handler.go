package inbound

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ahmad-banikamali/shortener/internal/core/domain"
	"github.com/ahmad-banikamali/shortener/internal/core/port"
)

type HTTPHandler struct {
	shorten port.ShortenUseCase
	resolve port.ResolveUseCase
	baseURL string
}

func NewHTTPHandler(s port.ShortenUseCase, r port.ResolveUseCase, baseURL string) *HTTPHandler {
	return &HTTPHandler{shorten: s, resolve: r, baseURL: baseURL}
}

func (h *HTTPHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /links", h.handleShorten)
	mux.HandleFunc("GET /{code}", h.handleResolve)
}

type shortenRequest struct {
	URL string `json:"url"`
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

	link, err := h.shorten.Shorten(r.Context(), req.URL)
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

// The adapter translates domain errors into HTTP status codes.
func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrInvalidURL):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, "internal error", http.StatusInternalServerError)
	}
}
