package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Expand handles URL expansion requests
func (h *ShortlinkHandler) Expand(c *gin.Context) {
	shortID := c.Param("shortID")
	if shortID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short ID"})
		return
	}

	// Call the injected URL service
	originalURL, err := h.URLService.ExpandURL(shortID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to expand URL"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
