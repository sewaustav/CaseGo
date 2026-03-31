package http

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware/rs256"
)

func (h *HttpHandler) RegisterRoutes(rg *gin.RouterGroup, middleware *rs256.JWTAuthMiddleware) {
	routes := rg.Group("/case_go")
	routes.Use(middleware.Handler())
	{
		routes.GET("/profile", h.GetProfileHandler)
		routes.GET("/history", h.GetHistoryHandler)
	}
}
