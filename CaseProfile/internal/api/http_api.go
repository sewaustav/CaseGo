package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/handlers/http"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware/rs256"
)

func SetupRoutes(httpHandler *http.HttpHandler, middleware *rs256.JWTAuthMiddleware) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // setup later
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")
	httpHandler.RegisterRoutes(v1, middleware)
	return r
}
