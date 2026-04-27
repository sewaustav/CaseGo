package profile_handler

import (
	"github.com/YoungFlores/Case_Go/Profile/pkg/middleware/rs256"
	"github.com/gin-gonic/gin"
)

func (h *ProfileHandler) RegisterRoutes(rg *gin.RouterGroup, jwtMiddleWare *rs256.JWTAuthMiddleware) {
	routers := rg.Group("/profile")
	routers.Use(jwtMiddleWare.Handler())
	{
		routers.POST("", h.CreateProfileHandler)
		routers.GET("", h.GetUserProfileHandler)
		routers.GET("/all", h.GetAllUsersHandler)
		routers.GET("/:id", h.GetUserByProfileIDHandler)
		routers.PUT("", h.UpdateProfileHandler)
		routers.PATCH("", h.PatchProfileHandler)
		routers.DELETE("", h.DeleteProfileHandler)
		routers.DELETE("/:id", h.HardDeleteHandler)
		
		routers.POST("/social", h.AddSocialLinkHandler)
		routers.PUT("/social/:id", h.UpdateLinkHandler)
		routers.DELETE("/social/:id", h.DeleteSocialLinkHandler)

		routers.POST("/purpose", h.AddPurposesHandler)
		routers.PUT("/purpose/:id", h.UpdatePurposeHandler)
		routers.DELETE("/purpose/:id", h.DeletePurposeHandler)

		routers.POST("/profession", h.AddProfessionsHandler)
		routers.GET("/profession", h.GetProfessionsHandler)
		routers.PUT("/profession/:id", h.EditProfessionsHandler)
		routers.DELETE("/profession/:id", h.DeleteProfessionsHandler)
	}
}
