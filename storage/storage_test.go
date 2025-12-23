package storage_test

import (
	"testing"

	"github.com/abdinurelmi/url-shortener/storage"
	"github.com/alicebob/miniredis/v2"
)

func TestStorage(t *testing.T) {
	mr, err := miniredis.Run()
	if err != nil {
		t.Fatalf("Failed to start miniredis: %v", err)
	}
	defer mr.Close()

	s, err := storage.NewService(mr.Addr())
	if err != nil {
		t.Fatalf("Failed to create storage service: %v", err)
	}

	t.Run("Save and Get URL", func(t *testing.T) {
		short := "abc"
		long := "https://example.com"

		if err := s.SaveURL(short, long); err != nil {
			t.Errorf("SaveURL failed: %v", err)
		}

		got, err := s.GetURL(short)
		if err != nil {
			t.Errorf("GetURL failed: %v", err)
		}
		if got != long {
			t.Errorf("Expected %s, got %s", long, got)
		}
	})

	t.Run("Analytics", func(t *testing.T) {
		code := "stats"

		// Initial check
		clicks, err := s.GetClicks(code)
		if err != nil {
			t.Errorf("GetClicks failed: %v", err)
		}
		if clicks != 0 {
			t.Errorf("Expected 0 clicks, got %d", clicks)
		}

		// Increment
		if err := s.IncrementClicks(code); err != nil {
			t.Errorf("IncrementClicks failed: %v", err)
		}

		// Check again
		clicks, err = s.GetClicks(code)
		if err != nil {
			t.Errorf("GetClicks failed: %v", err)
		}
		if clicks != 1 {
			t.Errorf("Expected 1 click, got %d", clicks)
		}
	})
}
