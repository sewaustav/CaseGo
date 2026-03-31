package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/service"
)

type HttpHandler struct {
	service service.Service
}

func NewHttpHandler(service service.Service) *HttpHandler {
	return &HttpHandler{service: service}
}

func (h *HttpHandler) GetProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exists := h.GetUserID(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   role,
	}

	profile, err := h.service.GetProfileService(ctx, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *HttpHandler) GetHistoryHandler(c *gin.Context) {
	ctx := c.Request.Context()

	fromDateStr := c.Query("from")
	if fromDateStr == "" {
		fromDateStr = "2026-01-01"
	}
	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, role, exists := h.GetUserID(c)
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user := models.UserIdentity{
		UserID: userID,
		Role:   role,
	}

	history, err := h.service.GetHistoryService(ctx, user, fromDate)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
