package http_handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	service "github.com/sewaustav/CaseGoCore/internal/cases/service/core"
)

type CaseGoHttpHandler struct {
	service service.CaseGoService
}

func NewCaseGoHttpHandler(service service.CaseGoService) *CaseGoHttpHandler {
	return &CaseGoHttpHandler{service: service}
}

func (h *CaseGoHttpHandler) GetCasesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var req *dto.GetCasesDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cases, err := h.service.GetCasesService(ctx, req.Limit, req.Page, &dto.UserSettingsDto{
		Topic:    req.Topic,
		Category: req.Category,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cases)

}

func (h *CaseGoHttpHandler) StartCaseHandler(c *gin.Context) {
	ctx := c.Request.Context()

	uid, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user := models.UserIdentity{
		UserID: uid,
		Role:   role,
	}

	caseID, err := strconv.ParseInt(c.Param("caseID"), 10, 64)

	dialog, err := h.service.StartDialogService(ctx, caseID, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, dialog)

}

func (h *CaseGoHttpHandler) DialogHandler(c *gin.Context) {
	ctx := c.Request.Context()

	uid, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user := models.UserIdentity{
		UserID: uid,
		Role:   role,
	}

	var req *dto.InteractionDto
	if err := c.ShouldBindBodyWithJSON(req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	question, err := h.service.HandleInteractionService(ctx, req, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, question)

}

func (h *CaseGoHttpHandler) CompleteDialogHandler(c *gin.Context) {
	ctx := c.Request.Context()

	uid, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user := models.UserIdentity{
		UserID: uid,
		Role:   role,
	}

	dialogID, err := strconv.ParseInt(c.Param("dialogID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid dialog ID"})
		return
	}

	result, err := h.service.CompleteDialogService(ctx, dialogID, user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, result)
}

func (h *CaseGoHttpHandler) GetCaseByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	caseID, err := strconv.ParseInt(c.Param("caseID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid case ID"})
		return
	}

	caseModel, err := h.service.GetCaseByIDService(ctx, caseID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, caseModel)
}

func (h *CaseGoHttpHandler) GetUsersDialogsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	uid, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID, err := strconv.ParseInt(c.Query("userID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid page"})
		return
	}

	user := models.UserIdentity{
		UserID: uid,
		Role:   role,
	}

	conv, err := h.service.GetUsersDialogsService(ctx, user, userID, limit, page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conv)
}

func (h *CaseGoHttpHandler) GetDialogByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	dialogID, err := strconv.ParseInt(c.Param("dialogID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid dialog ID"})
		return
	}

	uid, role, exist := h.GetUserID(c)
	if !exist {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user := models.UserIdentity{
		UserID: uid,
		Role:   role,
	}

	if dialogID <= 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid dialog ID"})
		return
	}

	dialog, err := h.service.GetUserDialogByIDService(ctx, user, dialogID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dialog)
}
