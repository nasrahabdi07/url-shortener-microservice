package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/abdinurelmi/url-shortener/shortener"
	"github.com/abdinurelmi/url-shortener/storage"
)

type Handler struct {
	storage *storage.Service
	baseURL string
}

func NewHandler(storage *storage.Service, baseURL string) *Handler {
	return &Handler{
		storage: storage,
		baseURL: baseURL,
	}
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL  string `json:"short_url"`
	ShortCode string `json:"short_code"`
}

type AnalyticsResponse struct {
	ShortCode string `json:"short_code"`
	Clicks    int64  `json:"clicks"`
}

// ShortenURL handles creating a new short URL
func (h *Handler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	code := shortener.GenerateShortCode()
	// Retry loop for collisions could go here, omitting for simplicity as per requirements

	if err := h.storage.SaveURL(code, req.URL); err != nil {
		http.Error(w, "Failed to save URL", http.StatusInternalServerError)
		return
	}

	resp := ShortenResponse{
		ShortURL:  h.baseURL + "/" + code,
		ShortCode: code,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Redirect handles the redirection and analytics tracking
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	// Path should be like /{code}
	code := r.URL.Path[1:] // strip leading slash
	if code == "" {
		http.NotFound(w, r)
		return
	}

	originalURL, err := h.storage.GetURL(code)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Increment clicks asynchronously to not block redirect
	go h.storage.IncrementClicks(code)

	http.Redirect(w, r, originalURL, http.StatusFound)
}

// GetAnalytics returns click stats
func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	// Path pattern /analytics/{code}
	// We need to extract code manually since we aren't using a fancy router library
	// E.g. Path: /analytics/abcde

	// Basic parsing:
	prefix := "/analytics/"
	if len(r.URL.Path) <= len(prefix) {
		http.NotFound(w, r)
		return
	}
	code := r.URL.Path[len(prefix):]

	clicks, err := h.storage.GetClicks(code)
	if err != nil {
		http.Error(w, "Failed to retrieve analytics", http.StatusInternalServerError)
		return
	}

	resp := AnalyticsResponse{
		ShortCode: code,
		Clicks:    clicks,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
