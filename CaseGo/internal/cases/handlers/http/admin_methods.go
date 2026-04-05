package http_handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (h *CaseGoHttpHandler) CreateCaseHandler(c *gin.Context) {
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

	var req *dto.NewCaseDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	newCase, err := h.service.CreateCaseService(ctx, req, user)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, newCase)
}

func (h *CaseGoHttpHandler) UpdateCaseHandler(c *gin.Context) {
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
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid case ID"})
		return
	}

	var req *dto.NewCaseDto
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if req.Topic == nil || req.Description == nil || req.Category == nil || req.FirstQuestion == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	if *req.Topic == "" || *req.Description == "" || *req.Category == 0 || *req.FirstQuestion == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}

	updatedCase, err := h.service.PatchCaseService(ctx, caseID, req, user)
	if err != nil {
		HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, updatedCase)
}

func (h *CaseGoHttpHandler) DeleteCaseHandler(c *gin.Context) {
	ctx := c.Request.Context()

	caseID, err := strconv.ParseInt(c.Param("caseID"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid case ID"})
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

	if err := h.service.DeleteCaseService(ctx, caseID, user); err != nil {
		HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
