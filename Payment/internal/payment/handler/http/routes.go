package http_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/Payment/pkg/middleware/rs256"
)

func (h *PaymentHttpHandler) RegisterRoutes(rg *gin.RouterGroup, middleware *rs256.JWTAuthMiddleware) {
	routes := rg.Group("/payment")
	routes.Use(middleware.Handler())
	{
		
	}
}