package http_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/Payment/internal/payment/models"
)

const (
	UserIDKey = "sub"
	RoleKey   = "role"
)

func (h *PaymentHttpHandler) GetUserID(c *gin.Context) (int64, models.UserRole, bool) {
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


