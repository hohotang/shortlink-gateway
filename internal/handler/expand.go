package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hohotang/shortlink-gateway/internal/middleware"
	"go.uber.org/zap"
)

// Expand handles URL expansion requests
// @Summary      Expand a short URL
// @Description  Redirects to the original URL from a short URL ID
// @Tags         urls
// @Produce      html
// @Param        shortID  path      string  true  "Short URL ID"
// @Success      302      {string}  string  "Redirect to original URL"
// @Failure      400      {object}  map[string]string  "Bad Request"
// @Failure      500      {object}  map[string]string  "Internal Server Error"
// @Router       /v1/expand/{shortID} [get]
func (h *ShortlinkHandler) Expand(c *gin.Context) {
	shortID := c.Param("shortID")
	if shortID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing short ID"})
		return
	}

	// Call the injected URL service with request context
	originalURL, err := h.URLService.ExpandURL(c.Request.Context(), shortID)
	if err != nil {
		logger := middleware.GetLogger(c.Request.Context())
		logger.Error("Failed to expand URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to expand URL"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}
