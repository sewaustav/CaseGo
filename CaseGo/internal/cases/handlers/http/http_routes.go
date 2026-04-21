package http_handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/pkg/middleware/rs256"
)

func (h *CaseGoHttpHandler) RegisterRoutes(rg *gin.RouterGroup, middleware *rs256.JWTAuthMiddleware) {
	routes := rg.Group("/case_go")
	routes.Use(middleware.Handler())
	{
		routes.POST("/cases/:caseID", h.StartCaseHandler)
		routes.GET("/cases", h.GetCasesHandler)
		routes.GET("/cases/:caseID", h.GetCaseByIDHandler)
		routes.GET("/users/:userID/dialogs", h.GetUsersDialogsHandler)
		routes.POST("/dialog", h.DialogHandler)
		routes.GET("/dialogs/:dialogID", h.GetDialogByIDHandler)
		routes.POST("/dialogs/:dialogID/complete", h.CompleteDialogHandler)

		routes.POST("/case", h.CreateCaseHandler)
		routes.PUT("/case/:caseID", h.UpdateCaseHandler)
		routes.DELETE("/case/:caseID", h.DeleteCaseHandler)
	}
}
