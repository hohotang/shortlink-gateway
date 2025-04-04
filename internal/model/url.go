package model

// ShortenRequest represents a request to shorten a URL
type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}
