package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Expand(c *gin.Context) {
	shortID := c.Param("shortID")
	if shortID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short ID"})
		return
	}

	// TODO: 呼叫 gRPC 的 url-service
	originalURL := "https://example.com/original-url"

	c.Redirect(http.StatusFound, originalURL)
}
