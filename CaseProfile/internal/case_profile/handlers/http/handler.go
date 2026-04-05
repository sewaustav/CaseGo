package http

import (
	"net/http"
	"strconv"
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

func (h HttpHandler) GetProfileHandler(c *gin.Context) {
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

func (h HttpHandler) GetHistoryHandler(c *gin.Context) {
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

func (h HttpHandler) GetUserProfileHandler(c *gin.Context) {
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

	userIdStr := c.Param("user_id")
	idStr := c.Param("id")

	if userIdStr == "" && idStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID or case ID must be provided"})
		return
	}

	var userProfile *models.CaseProfile
	if userIdStr != "" {
		userIDReq, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		userProfile, err = h.service.GetProfileByUserIDService(ctx, userIDReq, user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	} else {
		idReq, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid case ID"})
			return
		}
		userProfile, err = h.service.GetProfileByIDService(ctx, idReq, user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, userProfile)
}

func (h HttpHandler) GetUserProfileHistoryHandler(c *gin.Context) {
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

	userIdStr := c.Param("user_id")
	if userIdStr == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "User ID must be provided"})
		return
	}

	userIdReq, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	history, err := h.service.GetUserHistoryService(ctx, userIdReq, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}

func (h HttpHandler) DeleteResultByIDHandler(c *gin.Context) {
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

	userIdStr := c.Param("user_id")
	userIdReq, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := h.service.DeleteResultByIDService(ctx, userIdReq, user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Result deleted successfully"})
}
