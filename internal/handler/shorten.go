package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hohotang/shortlink-gateway/internal/service"
)

type ShortlinkHandler struct {
	// Dependencies can be injected here
	URLService service.URLService
}

// NewShortlinkHandler creates a new ShortlinkHandler with the given URLService
func NewShortlinkHandler(urlService service.URLService) *ShortlinkHandler {
	return &ShortlinkHandler{
		URLService: urlService,
	}
}

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

// Shorten handles URL shortening requests
func (h *ShortlinkHandler) Shorten(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.OriginalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Call the injected URL service
	shortID, err := h.URLService.ShortenURL(req.OriginalURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to shorten URL"})
		return
	}

	shortURL := "http://localhost:8080/" + shortID

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}
