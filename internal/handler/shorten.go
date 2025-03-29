package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ShortenRequest struct {
	OriginalURL string `json:"original_url"`
}

func Shorten(c *gin.Context) {
	var req ShortenRequest

	if err := c.ShouldBindJSON(&req); err != nil || req.OriginalURL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// TODO: 呼叫 gRPC 的 url-service
	shortID := "abc123"
	shortURL := "http://localhost:8080/" + shortID

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}
