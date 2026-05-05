package http_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/Payment/pkg/middleware/rs256"
)

func (h *PaymentHttpHandler) RegisterRoutes(rg *gin.RouterGroup, middleware *rs256.JWTAuthMiddleware) {
	routes := rg.Group("/payment")
	routes.Use(middleware.Handler())
	{
		routes.GET("/subscription", h.GetMySubscriptionInfoHandler)
				routes.GET("/history", h.GetMyPaymentsHandler)
				routes.PATCH("/subscription", h.UpdateSubscriptionHandler)

				routes.DELETE("/users/:user_id", h.DeleteUserHandler)
				routes.GET("/users/profile", h.GetUserProfileHandler)         
				routes.GET("/users/payments", h.GetUsersPaymentsHandler)      
				
				routes.GET("/transactions/:transaction_id", h.GetPaymentByTransactionIDHandler) 
				routes.GET("/:id", h.GetPaymentByIDHandler)
	}
}