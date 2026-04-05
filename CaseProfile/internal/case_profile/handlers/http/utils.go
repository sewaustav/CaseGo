package http

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/apperrors"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
)

const (
	UserIDKey = "sub"
	RoleKey   = "role"
)

func (h *HttpHandler) GetUserID(c *gin.Context) (int64, models.UserRole, bool) {
	userIDAny, exists := c.Get(UserIDKey)
	if !exists {
		return 0, 0, false
	}

	roleAny, exists := c.Get(RoleKey)
	if !exists {
		return 0, 0, false
	}

	uid, ok := userIDAny.(int64)
	if !ok {
		return 0, 0, false
	}

	role, ok := roleAny.(int)
	if !ok {
		return 0, 0, false
	}

	return uid, models.UserRole(role), true
}

func HandleError(c *gin.Context, err error) {
	if appErr, ok := errors.AsType[*apperrors.AppError](err); ok {
		if appErr.Code == http.StatusInternalServerError {
			log.Printf("[ERROR] %s: %v", appErr.Message, appErr.Err)
			c.AbortWithStatusJSON(appErr.Code, gin.H{"error": "internal server error"})
			return
		}
		c.AbortWithStatusJSON(appErr.Code, gin.H{"error": appErr.Message})
		return
	}

	log.Printf("[ERROR] unhandled error: %v", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
