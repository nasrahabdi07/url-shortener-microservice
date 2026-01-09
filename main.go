package main

import (
	"log"
	"net/http"
	"os"

	"github.com/abdinurelmi/url-shortener/handlers"
	"github.com/abdinurelmi/url-shortener/storage"
)

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}

	log.Printf("Connecting to Redis at %s...", redisAddr)
	store, err := storage.NewService(redisAddr)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	handler := handlers.NewHandler(store, baseURL)

	// Using standard http.ServeMux
	mux := http.NewServeMux()

	// Register routes
	mux.HandleFunc("/shorten", handler.ShortenURL)
	mux.HandleFunc("/analytics/", handler.GetAnalytics)
	// Catch-all for redirects - must be registered as "/" to match all paths,
	// but we need to be careful not to override specific paths if we used them differently.
	// Since /shorten and /analytics/ are specific, "/" will catch everything else.
	// However, http.ServeMux "longest match" rule applies.
	mux.HandleFunc("/", handler.Redirect)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	log.Printf("Starting server on :8080 with base URL %s", baseURL)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
