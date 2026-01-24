package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/pkg/middleware"
)

func SetupRouter(handler *ProfileHandler, jwtMiddleware *middleware.JWTAuthMiddleware) *gin.Engine {
	r := gin.Default()
	
	protected := r.Group("profile/api/v1")
	protected.Use()
	{
		protected.POST("/profile", handler.CreateProfileHandler)
		protected.GET("/profile", handler.GetUserProfileHandler)
		protected.PUT("/profile", handler.UpdateProfileHandler)
		protected.PATCH("/profile", handler.PatchProfileHandler)
		protected.DELETE("/profile", handler.DeleteProfileHandler)
		protected.DELETE("/profile/:id", handler.HardDeleteHandler)

		protected.POST("/profile/social", handler.AddSocialLinkHandler)
		protected.PUT("/profile/social/:id", handler.UpdateLinkHandler)
		protected.DELETE("profile/social/:id", handler.DeleteSocialLinkHandler)

		protected.POST("/profile/purpose", handler.AddPurposesHandler)
		protected.PUT("/profile/purpose/:id", handler.UpdatePurposeHandler)
		protected.DELETE("/profile/purpose/:id", handler.DeletePuposeHandler)
	}


	return r
}