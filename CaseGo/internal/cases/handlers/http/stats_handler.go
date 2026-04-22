package http_handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (h *CaseGoHttpHandler) GetStatsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	_, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if role != models.Admin {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
		return
	}

	stats, err := h.service.GetStatsService(ctx)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}
