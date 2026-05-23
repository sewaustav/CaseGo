package api

import (
	_ "github.com/YoungFlores/Case_Go/Profile/docs"
	categoriesHandler "github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/handlers"
	profileHandler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers/http"
	"github.com/YoungFlores/Case_Go/Profile/internal/search/handlers"
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func SetupRouter(profileHandler *profileHandler.ProfileHandler, searchHandler *handlers.SearchHandler, categoryHandler *categoriesHandler.ProfessionCategoryHandler, jwtMiddleware *rs256.JWTAuthMiddleware) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // setup later
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/profile/api/v1")

	profileHandler.RegisterRoutes(v1, jwtMiddleware)
	searchHandler.RegisterRoutes(v1, jwtMiddleware)
	categoryHandler.RegisterRoutes(v1, jwtMiddleware)

	return r
}
