package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/internal/profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
	profileService "github.com/sewaustav/CaseGoProfile/internal/profile/service"
)

type ProfileHandler struct {
	service *profileService.ProfileService
}

func NewProfileHandler(service *profileService.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

// POST
func (h *ProfileHandler) CreateProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body dto.CreateProfileRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.service.CreateProfileService(ctx, body, userInfo)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, profile)

}

func (h *ProfileHandler) AddSocialLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body []dto.SocialLinkDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	links, err := h.service.AddSocialLinksService(ctx, body, userInfo)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, links)
}

func (h *ProfileHandler) AddPurposesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body []dto.UserPurposeDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	purposes, err := h.service.AddPurposesService(ctx, body, userInfo)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusCreated, purposes)
	
}

// PATCH
func (h *ProfileHandler) PatchProfileHandler(c  *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body dto.UpdateProfilePartialDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.service.PatchProfileService(ctx, userInfo, body) 
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// PUT 
func (h *ProfileHandler) UpdateProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body dto.ProfileInfoDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.service.UpdateProfileService(ctx, userInfo, body)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}
	
	c.JSON(http.StatusOK, profile)
	
}

func (h *ProfileHandler) UpdateLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body dto.SocialLinkDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	links, err := h.service.UpdateSocialLinkService(ctx, body, userInfo, id)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusOK, links)
}

func (h *ProfileHandler) UpdatePurposeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	var body dto.UserPurposeDTO

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	purposes, err := h.service.UpdatePurposeService(ctx, body, userInfo, id)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusOK, purposes)
}

// Get 
func (h *ProfileHandler) GetUserProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	profile, err := h.service.GetUserProfileService(ctx, userInfo)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

func (h *ProfileHandler) GetUserByProfileIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	profile, err := h.service.GetUserProfileByIDService(ctx, userInfo, id)
	if err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.JSON(http.StatusOK, profile)
	
}


// Delete
func (h *ProfileHandler) DeleteProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	if err := h.service.DeleteProfileService(ctx, userInfo); err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *ProfileHandler) HardDeleteHandler(c *gin.Context) {
	ctx := c.Request.Context()

	usrID, err := strconv.ParseInt(c.Param("id"), 10, 64) // usrID - id of user who we want to delete
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	if err := h.service.DeleteProfileWithoutRecoveryService(ctx, userInfo, usrID); err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.Status(http.StatusNoContent)

}

func (h *ProfileHandler) DeletePuposeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	if err := h.service.DeletePuposeService(ctx, id, userInfo); err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.Status(http.StatusNoContent) 
	

}

func (h *ProfileHandler) DeleteSocialLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.getUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role 
	userInfo.UserID = userID

	if err := h.service.DeleteLinkService(ctx, id, userInfo); err != nil {
		// validate error in future
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create event"})
		return
	}

	c.Status(http.StatusNoContent)
}