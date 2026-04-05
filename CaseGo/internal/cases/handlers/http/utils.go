package http_handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

const (
	UserIDKey = "sub"
	RoleKey   = "role"
)

func (h *CaseGoHttpHandler) GetUserID(c *gin.Context) (int64, models.UserRole, bool) {
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

// HandleError маппит *apperrors.AppError на правильный HTTP-статус.
// Для 500 логируем подробности, клиенту отдаём общее сообщение.
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

	// Неизвестная ошибка — всегда 500, детали не светим
	log.Printf("[ERROR] unhandled error: %v", err)
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
