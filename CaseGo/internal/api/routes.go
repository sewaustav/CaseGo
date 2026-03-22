package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/internal/cases/handlers/http"
	"github.com/sewaustav/CaseGoCore/pkg/middleware/rs256"
)

func SetupRoutes(casesHandler *http_handlers.CaseGoHttpHandler, jwtMiddleware *rs256.JWTAuthMiddleware) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // setup later
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	v1 := r.Group("/api/v1")

	casesHandler.RegisterRoutes(v1, jwtMiddleware)
	return r
}
