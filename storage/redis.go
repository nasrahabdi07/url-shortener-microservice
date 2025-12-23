package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Service struct {
	client *redis.Client
}

func NewService(addr string) (*Service, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Service{client: client}, nil
}

// SaveURL stores the original URL with a short code (persistent for now)
func (s *Service) SaveURL(shortCode, originalURL string) error {
	// Key: url:{shortCode} -> originalURL
	key := fmt.Sprintf("url:%s", shortCode)
	// Using 0 expiration for persistent storage, or could add TTL
	return s.client.Set(ctx, key, originalURL, 0).Err()
}

// GetURL retrieves the original URL given a short code
func (s *Service) GetURL(shortCode string) (string, error) {
	key := fmt.Sprintf("url:%s", shortCode)
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("URL not found")
	}
	return val, err
}

// IncrementClicks increments the click count for a short code
func (s *Service) IncrementClicks(shortCode string) error {
	key := fmt.Sprintf("analytics:%s", shortCode)
	return s.client.Incr(ctx, key).Err()
}

// GetClicks retrieves the click count for a short code
func (s *Service) GetClicks(shortCode string) (int64, error) {
	key := fmt.Sprintf("analytics:%s", shortCode)
	val, err := s.client.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil // treat missing key as 0 clicks
	}
	return val, err
}
