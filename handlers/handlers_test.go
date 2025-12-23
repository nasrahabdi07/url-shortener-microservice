package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/abdinurelmi/url-shortener/handlers"
	"github.com/abdinurelmi/url-shortener/storage"
	"github.com/alicebob/miniredis/v2"
)

func TestHandlers(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	store, err := storage.NewService(mr.Addr())
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	h := handlers.NewHandler(store, "http://test.com")

	t.Run("Shorten URL", func(t *testing.T) {
		reqBody := `{"url": "https://google.com"}`
		req := httptest.NewRequest("POST", "/shorten", bytes.NewBufferString(reqBody))
		w := httptest.NewRecorder()

		h.ShortenURL(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		var res handlers.ShortenResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if res.ShortCode == "" {
			t.Error("Expected short_code, got empty")
		}
		if res.ShortURL == "" {
			t.Error("Expected short_url, got empty")
		}
	})

	t.Run("Redirect", func(t *testing.T) {
		// Setup data
		code := "xyz"
		long := "https://example.com"
		store.SaveURL(code, long)

		req := httptest.NewRequest("GET", "/"+code, nil)
		w := httptest.NewRecorder()

		h.Redirect(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusFound {
			t.Errorf("Expected 302, got %d", resp.StatusCode)
		}
		loc, _ := resp.Location()
		if loc.String() != long {
			t.Errorf("Expected redirect to %s, got %s", long, loc.String())
		}

		// Wait a bit for async increment
		time.Sleep(100 * time.Millisecond)

		clicks, _ := store.GetClicks(code)
		if clicks != 1 {
			t.Errorf("Expected 1 click, got %d", clicks)
		}
	})

	t.Run("Analytics", func(t *testing.T) {
		code := "ana"
		store.IncrementClicks(code)
		store.IncrementClicks(code)

		req := httptest.NewRequest("GET", "/analytics/"+code, nil)
		w := httptest.NewRecorder()

		h.GetAnalytics(w, req)

		resp := w.Result()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected 200, got %d", resp.StatusCode)
		}

		var res handlers.AnalyticsResponse
		if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if res.Clicks != 2 {
			t.Errorf("Expected 2 clicks, got %d", res.Clicks)
		}
	})
}
