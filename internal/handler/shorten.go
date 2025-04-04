package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hohotang/shortlink-gateway/internal/model"
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

// Shorten handles URL shortening requests
// @Summary      Shorten a URL
// @Description  Creates a short URL from a long URL
// @Tags         urls
// @Accept       json
// @Produce      json
// @Param        request  body      model.ShortenRequest  true  "URL to shorten"
// @Success      200      {object}  map[string]string  "Returns shortened URL"
// @Failure      400      {object}  map[string]string  "Bad Request"
// @Failure      500      {object}  map[string]string  "Internal Server Error"
// @Router       /shorten [post]
func (h *ShortlinkHandler) Shorten(c *gin.Context) {
	var req model.ShortenRequest

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
