package api

import (
	"github.com/gin-gonic/gin"
	http_handler "github.com/sewaustav/Payment/internal/payment/handler/http"
	"github.com/sewaustav/Payment/pkg/middleware/rs256"
	"github.com/gin-contrib/cors"
)

func SetupRoutes(handler *http_handler.PaymentHttpHandler, jwtMiddleware rs256.JWTAuthMiddleware) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // setup later
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")

	handler.RegisterRoutes(v1, &jwtMiddleware)

	return r
}