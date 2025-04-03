package service

import (
	"errors"
)

// URLService provides URL shortening functionality
type URLService interface {
	ShortenURL(originalURL string) (string, error)
	ExpandURL(shortID string) (string, error)
}

// URLServiceImpl implements URLService interface
type URLServiceImpl struct {
	client URLService // This can be a real gRPC client or a mock
}

// NewURLService creates a new URL service
func NewURLService() URLService {
	return &URLServiceImpl{
		client: &MockURLService{}, // Default to mock implementation
	}
}

// NewURLServiceWithClient creates a URL service with a specific client
func NewURLServiceWithClient(client URLService) URLService {
	return &URLServiceImpl{
		client: client,
	}
}

// ShortenURL creates a short URL from the original URL
func (s *URLServiceImpl) ShortenURL(originalURL string) (string, error) {
	return s.client.ShortenURL(originalURL)
}

// ExpandURL resolves a short URL to its original URL
func (s *URLServiceImpl) ExpandURL(shortID string) (string, error) {
	return s.client.ExpandURL(shortID)
}

// MockURLService provides a local implementation for testing/development
type MockURLService struct{}

// ShortenURL creates a short URL from the original URL
func (s *MockURLService) ShortenURL(originalURL string) (string, error) {
	// This is a mock implementation for testing
	if originalURL == "" {
		return "", errors.New("original URL cannot be empty")
	}

	return "abc123", nil
}

// ExpandURL resolves a short URL to its original URL
func (s *MockURLService) ExpandURL(shortID string) (string, error) {
	// This is a mock implementation for testing
	if shortID == "" {
		return "", errors.New("short ID cannot be empty")
	}

	return "https://example.com/original-url", nil
}
