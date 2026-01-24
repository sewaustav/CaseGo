package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
)

func (h *ProfileHandler) getUserID(c *gin.Context) (int64, models.UserRole, bool) {
	return 123, models.Guest, true
}