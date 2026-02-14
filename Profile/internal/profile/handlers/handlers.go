package profile_handler

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
	profileService "github.com/YoungFlores/Case_Go/Profile/internal/profile/service"
	apperrors "github.com/YoungFlores/Case_Go/Profile/pkg/errors"
	"github.com/gin-gonic/gin"
)

type ProfileHandler struct {
	service *profileService.ProfileService
}

func NewProfileHandler(service *profileService.ProfileService) *ProfileHandler {
	return &ProfileHandler{
		service: service,
	}
}

// CreateProfileHandler godoc
// @Summary Создать новый профиль
// @Description Создает профиль пользователя с соцсетями и целями
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.CreateProfileRequest true "Данные профиля"
// @Success 201 {object} models.Profile "Профиль успешно создан"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 409 {object} map[string]string "Conflict - Field already taken"
// @Router /profile [post]
func (h *ProfileHandler) CreateProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()
	slog.Info("Create profile request")
	userID, role, exist := h.GetUserID(c)
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
		var conflictErr *repoerr.RepoError

		if errors.As(err, &conflictErr) {
			resp := map[string]string{
				"error":   "Conflict",
				"field":   conflictErr.Field,
				"message": fmt.Sprintf("Value in field %s is already taken", conflictErr.Field),
			}
			c.JSON(http.StatusConflict, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, profile)

}

// AddSocialLinkHandler godoc
// @Summary Добавить ссылки на соцсети
// @Tags social
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body []dto.SocialLinkDTO true "Список ссылок"
// @Success 201 {array} dto.SocialLinkDTO
// @Router /profile/social [post]
func (h *ProfileHandler) AddSocialLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add link"})
		return
	}

	c.JSON(http.StatusCreated, links)
}

// AddPurposesHandler godoc
// @Summary Создать новые цели
// @Tags purpose
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body []dto.UserPurposeDTO true "Список целей"
// @Success 201 {array} dto.UserPurposeDTO
// @Router /profile/purposes [post]
func (h *ProfileHandler) AddPurposesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add purpose"})
		return
	}

	c.JSON(http.StatusCreated, purposes)

}

// PATCH

// PatchProfileHandler godoc
// @Summary Частично обновить профиль
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.UpdateProfilePartialDTO true "Поля для обновления"
// @Success 200 {object} models.Profile
// @Router /profile [patch]
func (h *ProfileHandler) PatchProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
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
		var conflictErr *repoerr.RepoError

		if errors.As(err, &conflictErr) {
			resp := map[string]string{
				"error":   "Conflict",
				"field":   conflictErr.Field,
				"message": fmt.Sprintf("Value in field %s is already taken", conflictErr.Field),
			}
			c.JSON(http.StatusConflict, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to patch profile"})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// PUT

// UpdateProfileHandler godoc
// @Summary Полное обновление профиля
// @Tags profile
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body dto.ProfileInfoDTO true "Новые данные профиля"
// @Success 200 {object} models.Profile
// @Failure 409 {object} map[string]string "Conflict - Field already taken"
// @Router /profile [put]
func (h *ProfileHandler) UpdateProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
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
		var conflictErr *repoerr.RepoError

		if errors.As(err, &conflictErr) {
			resp := map[string]string{
				"error":   "Conflict",
				"field":   conflictErr.Field,
				"message": fmt.Sprintf("Value in field %s is already taken", conflictErr.Field),
			}
			c.JSON(http.StatusConflict, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, profile)

}

// UpdateLinkHandler godoc
// @Summary Обновить конкретную соц. ссылку
// @Tags social
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID ссылки"
// @Param body body dto.SocialLinkDTO true "Данные ссылки"
// @Success 200 {object} dto.SocialLinkDTO
// @Failure 403 {object} map[string]string "Forbidden - Not your link"
// @Router /profile/social/{id} [put]
func (h *ProfileHandler) UpdateLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.GetUserID(c)
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
		var forbErr *repoerr.RepoError

		if errors.As(err, &forbErr) {
			resp := map[string]string{
				"error":   "Forbidden",
				"field":   forbErr.Field,
				"message": fmt.Sprintf("Edit id %s", forbErr.Field),
			}
			c.JSON(http.StatusForbidden, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update link"})
		return
	}

	c.JSON(http.StatusOK, links)
}

// @Summary Обновить конкретную цель
// @Tags purpose
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID цели"
// @Param body body dto.UserPurposeDTO true "Данные цели"
// @Success 200 {object} dto.UserPurposeDTO
// @Failure 400 {object} map[string]string "Invalid ID"
// @Failure 403 {object} map[string]string "Forbidden"
// @Router /profiles/purpose/{id} [put]
func (h *ProfileHandler) UpdatePurposeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.GetUserID(c)
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
		var forbErr *repoerr.RepoError

		if errors.As(err, &forbErr) {
			resp := map[string]string{
				"error":   "Forbidden",
				"field":   forbErr.Field,
				"message": fmt.Sprintf("Edit id %s", forbErr.Field),
			}
			c.JSON(http.StatusForbidden, resp)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update purpose"})
		return
	}

	c.JSON(http.StatusOK, purposes)
}

// Get

// GetUserProfileHandler godoc
// @Summary Получить свой профиль
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Profile
// @Success 204 "Профиль не активен"
// @Router /profile [get]
func (h *ProfileHandler) GetUserProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role
	userInfo.UserID = userID

	profile, err := h.service.GetUserProfileService(ctx, userInfo)
	if err != nil {
		if errors.Is(err, apperrors.ErrIsNotActive) {
			c.JSON(http.StatusNotFound, gin.H{"info": "user is not active"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed %s", err)})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// GetUserByProfileIDHandler godoc
// @Summary Получить профиль по ID
// @Description Позволяет просмотреть чужой профиль по его ID
// @Tags profile
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID профиля"
// @Success 200 {object} models.Profile
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Router /profile/{id} [get]
func (h *ProfileHandler) GetUserByProfileIDHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role
	userInfo.UserID = userID

	profile, err := h.service.GetUserProfileByIDService(ctx, userInfo, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrForbidden) {
			c.Status(http.StatusForbidden)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.JSON(http.StatusOK, profile)

}

// Delete

// DeleteProfileHandler godoc
// @Summary Мягкое удаление своего профиля
// @Tags profile
// @Security BearerAuth
// @Success 204 "No Content"
// @Router /profile [delete]
func (h *ProfileHandler) DeleteProfileHandler(c *gin.Context) {
	ctx := c.Request.Context()

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role
	userInfo.UserID = userID

	if err := h.service.DeleteProfileService(ctx, userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.Status(http.StatusNoContent)

}

// HardDeleteHandler godoc
// @Summary Полное удаление профиля (Admin only)
// @Tags profile
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 204 "No Content"
// @Failure 403 "Forbidden"
// @Router /profile/{id} [delete]
func (h *ProfileHandler) HardDeleteHandler(c *gin.Context) {
	ctx := c.Request.Context()

	usrID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role
	userInfo.UserID = userID

	if err := h.service.DeleteProfileWithoutRecoveryService(ctx, userInfo, usrID); err != nil {
		if errors.Is(err, apperrors.ErrForbidden) {
			c.Status(http.StatusForbidden)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.Status(http.StatusNoContent)

}

// DeletePuposeHandler godoc
// @Summary Удалить цель
// @Tags purpose
// @Security BearerAuth
// @Param id path int true "ID цели"
// @Success 204 "No Content"
// @Failure 403 "Forbidden"
// @Router /profile/purpose/{id} [delete]
func (h *ProfileHandler) DeletePurposeHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role
	userInfo.UserID = userID

	if err := h.service.DeletePurposeService(ctx, id, userInfo); err != nil {
		if errors.Is(err, apperrors.ErrForbidden) {
			c.Status(http.StatusForbidden)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed "})
		return
	}

	c.Status(http.StatusNoContent)

}

// DeleteSocialLinkHandler godoc
// @Summary Удалить соц. ссылку
// @Tags social
// @Security BearerAuth
// @Param id path int true "ID ссылки"
// @Success 204 "No Content"
// @Failure 403 "Forbidden"
// @Router /profile/social/{id} [delete]
func (h *ProfileHandler) DeleteSocialLinkHandler(c *gin.Context) {
	ctx := c.Request.Context()

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	userID, role, exist := h.GetUserID(c)
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var userInfo models.UserIdentity
	userInfo.Role = role
	userInfo.UserID = userID

	if err := h.service.DeleteLinkService(ctx, id, userInfo); err != nil {
		if errors.Is(err, apperrors.ErrForbidden) {
			c.Status(http.StatusForbidden)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
