package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	profile_handler "github.com/YoungFlores/Case_Go/Profile/internal/profile/handlers/http"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
	"github.com/YoungFlores/Case_Go/Profile/mocks"
	apperrors "github.com/YoungFlores/Case_Go/Profile/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func ptrString(s string) *string { return &s }
func ptrInt(i int) *int          { return &i }

func setupRouter(mockSvc *mocks.ProfileCore, handler func(*gin.Context), method, path string, userID int64, role models.UserRole) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Handle(method, path, func(c *gin.Context) {
		c.Set("sub", userID)
		c.Set("role", int(role))
		c.Next()
	}, handler)
	return r
}

func setupRouterNoAuth(handler func(*gin.Context), method, path string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Handle(method, path, handler)
	return r
}

// ==================== CreateProfileHandler Tests ====================

func TestCreateProfileHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)
		userID := int64(1)

		reqBody := dto.CreateProfileRequest{
			Info: dto.ProfileInfoDTO{
				Avatar:     "https://avatar.com",
				Name:       "Маша",
				Surname:    "Залужная",
				Username:   "masha_zaluzhnaya",
				City:       ptrString("Moscow"),
				Age:        ptrInt(21),
				Sex:        ptrInt(1),
				Profession: ptrString("Проектировщик ракет"),
			},
			SocialLinks: []dto.SocialLinkDTO{
				{Type: "telegram", URL: "https://t.me/MashaZalushnaya"},
			},
			Purposes: []dto.UserPurposeDTO{
				{Purpose: "Донбас"},
			},
		}

		expectedProfile := &models.Profile{
			ID:       1,
			UserID:   userID,
			Avatar:   "https://avatar.com",
			IsActive: true,
			Username: "masha_zaluzhnaya",
			Name:     "Маша",
			Surname:  "Залужная",
		}

		expectedResponse := &models.UserProfile{
			UsrProfile:  *expectedProfile,
			UsrPurposes: []models.UserPurpose{{ID: 1, Purpose: "Донбас", UserID: userID}},
			UsrSocials:  []models.UserSocialLink{{ID: 1, Type: "telegram", URL: "https://t.me/MashaZalushnaya", UserID: userID}},
		}

		mockSvc.On("CreateProfileService", mock.Anything, reqBody, mock.Anything).
			Return(expectedResponse, nil).Once()

		r := setupRouter(mockSvc, h.CreateProfileHandler, http.MethodPost, "/profiles", userID, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.CreateProfileHandler, http.MethodPost, "/profiles")

		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.CreateProfileHandler, http.MethodPost, "/profiles", 1, models.User)

		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewReader([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("conflict error - username taken", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.CreateProfileRequest{
			Info: dto.ProfileInfoDTO{
				Avatar:   "https://avatar.com",
				Name:     "Маша",
				Surname:  "Залужная",
				Username: "existing_username",
			},
			Purposes: []dto.UserPurposeDTO{{Purpose: "Test purpose"}},
		}

		mockSvc.On("CreateProfileService", mock.Anything, reqBody, mock.Anything).
			Return(nil, &repoerr.RepoError{Field: "username", Err: repoerr.ErrConflict}).Once()

		r := setupRouter(mockSvc, h.CreateProfileHandler, http.MethodPost, "/profiles", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.CreateProfileRequest{
			Info: dto.ProfileInfoDTO{
				Avatar:   "https://avatar.com",
				Name:     "Маша",
				Surname:  "Залужная",
				Username: "masha",
			},
			Purposes: []dto.UserPurposeDTO{{Purpose: "Test purpose"}},
		}

		mockSvc.On("CreateProfileService", mock.Anything, reqBody, mock.Anything).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.CreateProfileHandler, http.MethodPost, "/profiles", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profiles", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== GetUserProfileHandler Tests ====================

func TestGetUserProfileHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)
		userID := int64(1)

		expectedResponse := &models.UserProfile{
			UsrProfile: models.Profile{
				ID:       1,
				UserID:   userID,
				Name:     "Test",
				IsActive: true,
			},
			UsrPurposes: []models.UserPurpose{},
			UsrSocials:  []models.UserSocialLink{},
		}

		mockSvc.On("GetUserProfileService", mock.Anything, mock.Anything).
			Return(expectedResponse, nil).Once()

		r := setupRouter(mockSvc, h.GetUserProfileHandler, http.MethodGet, "/profile", userID, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.GetUserProfileHandler, http.MethodGet, "/profile")

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("user not active", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("GetUserProfileService", mock.Anything, mock.Anything).
			Return(nil, apperrors.ErrIsNotActive).Once()

		r := setupRouter(mockSvc, h.GetUserProfileHandler, http.MethodGet, "/profile", 1, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("GetUserProfileService", mock.Anything, mock.Anything).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.GetUserProfileHandler, http.MethodGet, "/profile", 1, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== GetUserByProfileIDHandler Tests ====================

func TestGetUserByProfileIDHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		expectedResponse := &models.UserProfile{
			UsrProfile: models.Profile{ID: 5, UserID: 1, Name: "Test"},
		}

		mockSvc.On("GetUserProfileByIDService", mock.Anything, mock.Anything, int64(5)).
			Return(expectedResponse, nil).Once()

		r := setupRouter(mockSvc, h.GetUserByProfileIDHandler, http.MethodGet, "/profile/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.GetUserByProfileIDHandler, http.MethodGet, "/profile/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.GetUserByProfileIDHandler, http.MethodGet, "/profile/:id")

		req, _ := http.NewRequest(http.MethodGet, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("GetUserProfileByIDService", mock.Anything, mock.Anything, int64(5)).
			Return(nil, apperrors.ErrForbidden).Once()

		r := setupRouter(mockSvc, h.GetUserByProfileIDHandler, http.MethodGet, "/profile/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("GetUserProfileByIDService", mock.Anything, mock.Anything, int64(5)).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.GetUserByProfileIDHandler, http.MethodGet, "/profile/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodGet, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== UpdateProfileHandler Tests ====================

func TestUpdateProfileHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.ProfileInfoDTO{
			Avatar:   "https://new-avatar.com",
			Name:     "NewName",
			Surname:  "NewSurname",
			Username: "newusername",
		}

		expectedProfile := &models.Profile{
			ID:       1,
			UserID:   1,
			Avatar:   "https://new-avatar.com",
			Name:     "NewName",
			Username: "newusername",
		}

		mockSvc.On("UpdateProfileService", mock.Anything, mock.Anything, reqBody).
			Return(expectedProfile, nil).Once()

		r := setupRouter(mockSvc, h.UpdateProfileHandler, http.MethodPut, "/profile", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.UpdateProfileHandler, http.MethodPut, "/profile")

		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.UpdateProfileHandler, http.MethodPut, "/profile", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("conflict error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.ProfileInfoDTO{
			Avatar:   "https://avatar.com",
			Name:     "Name",
			Surname:  "Surname",
			Username: "taken_username",
		}

		mockSvc.On("UpdateProfileService", mock.Anything, mock.Anything, reqBody).
			Return(nil, &repoerr.RepoError{Field: "username", Err: repoerr.ErrConflict}).Once()

		r := setupRouter(mockSvc, h.UpdateProfileHandler, http.MethodPut, "/profile", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.ProfileInfoDTO{
			Avatar:   "https://avatar.com",
			Name:     "Name",
			Surname:  "Surname",
			Username: "username",
		}

		mockSvc.On("UpdateProfileService", mock.Anything, mock.Anything, reqBody).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.UpdateProfileHandler, http.MethodPut, "/profile", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== PatchProfileHandler Tests ====================

func TestPatchProfileHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		newName := "UpdatedName"
		reqBody := dto.UpdateProfilePartialDTO{
			Name: &newName,
		}

		expectedProfile := &models.Profile{
			ID:     1,
			UserID: 1,
			Name:   newName,
		}

		mockSvc.On("PatchProfileService", mock.Anything, mock.Anything, reqBody).
			Return(expectedProfile, nil).Once()

		r := setupRouter(mockSvc, h.PatchProfileHandler, http.MethodPatch, "/profile", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, "/profile", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.PatchProfileHandler, http.MethodPatch, "/profile")

		req, _ := http.NewRequest(http.MethodPatch, "/profile", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.PatchProfileHandler, http.MethodPatch, "/profile", 1, models.User)

		req, _ := http.NewRequest(http.MethodPatch, "/profile", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("conflict error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		takenUsername := "taken"
		reqBody := dto.UpdateProfilePartialDTO{
			Username: &takenUsername,
		}

		mockSvc.On("PatchProfileService", mock.Anything, mock.Anything, reqBody).
			Return(nil, &repoerr.RepoError{Field: "username", Err: repoerr.ErrConflict}).Once()

		r := setupRouter(mockSvc, h.PatchProfileHandler, http.MethodPatch, "/profile", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, "/profile", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		name := "name"
		reqBody := dto.UpdateProfilePartialDTO{
			Name: &name,
		}

		mockSvc.On("PatchProfileService", mock.Anything, mock.Anything, reqBody).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.PatchProfileHandler, http.MethodPatch, "/profile", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPatch, "/profile", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== DeleteProfileHandler Tests ====================

func TestDeleteProfileHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfileService", mock.Anything, mock.Anything).
			Return(nil).Once()

		r := setupRouter(mockSvc, h.DeleteProfileHandler, http.MethodDelete, "/profile", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.DeleteProfileHandler, http.MethodDelete, "/profile")

		req, _ := http.NewRequest(http.MethodDelete, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfileService", mock.Anything, mock.Anything).
			Return(errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.DeleteProfileHandler, http.MethodDelete, "/profile", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== HardDeleteHandler Tests ====================

func TestHardDeleteHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfileWithoutRecoveryService", mock.Anything, mock.Anything, int64(5)).
			Return(nil).Once()

		r := setupRouter(mockSvc, h.HardDeleteHandler, http.MethodDelete, "/profile/:id", 1, models.Admin)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.HardDeleteHandler, http.MethodDelete, "/profile/:id", 1, models.Admin)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.HardDeleteHandler, http.MethodDelete, "/profile/:id")

		req, _ := http.NewRequest(http.MethodDelete, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden - not admin", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfileWithoutRecoveryService", mock.Anything, mock.Anything, int64(5)).
			Return(apperrors.ErrForbidden).Once()

		r := setupRouter(mockSvc, h.HardDeleteHandler, http.MethodDelete, "/profile/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfileWithoutRecoveryService", mock.Anything, mock.Anything, int64(5)).
			Return(errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.HardDeleteHandler, http.MethodDelete, "/profile/:id", 1, models.Admin)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== AddSocialLinkHandler Tests ====================

func TestAddSocialLinkHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := []dto.SocialLinkDTO{
			{Type: "telegram", URL: "https://t.me/test"},
			{Type: "github", URL: "https://github.com/test"},
		}

		expectedLinks := []models.UserSocialLink{
			{ID: 1, UserID: 1, Type: "telegram", URL: "https://t.me/test"},
			{ID: 2, UserID: 1, Type: "github", URL: "https://github.com/test"},
		}

		mockSvc.On("AddSocialLinksService", mock.Anything, reqBody, mock.Anything).
			Return(expectedLinks, nil).Once()

		r := setupRouter(mockSvc, h.AddSocialLinkHandler, http.MethodPost, "/profile/social", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profile/social", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.AddSocialLinkHandler, http.MethodPost, "/profile/social")

		req, _ := http.NewRequest(http.MethodPost, "/profile/social", bytes.NewReader([]byte("[]")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.AddSocialLinkHandler, http.MethodPost, "/profile/social", 1, models.User)

		req, _ := http.NewRequest(http.MethodPost, "/profile/social", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := []dto.SocialLinkDTO{
			{Type: "telegram", URL: "https://t.me/test"},
		}

		mockSvc.On("AddSocialLinksService", mock.Anything, reqBody, mock.Anything).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.AddSocialLinkHandler, http.MethodPost, "/profile/social", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profile/social", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== UpdateLinkHandler Tests ====================

func TestUpdateLinkHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.SocialLinkDTO{Type: "twitter", URL: "https://twitter.com/test"}

		expectedLinks := []models.UserSocialLink{
			{ID: 5, UserID: 1, Type: "twitter", URL: "https://twitter.com/test"},
		}

		mockSvc.On("UpdateSocialLinkService", mock.Anything, reqBody, mock.Anything, int64(5)).
			Return(expectedLinks, nil).Once()

		r := setupRouter(mockSvc, h.UpdateLinkHandler, http.MethodPut, "/profile/social/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/social/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.UpdateLinkHandler, http.MethodPut, "/profile/social/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile/social/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.UpdateLinkHandler, http.MethodPut, "/profile/social/:id")

		req, _ := http.NewRequest(http.MethodPut, "/profile/social/5", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.UpdateLinkHandler, http.MethodPut, "/profile/social/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile/social/5", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.SocialLinkDTO{Type: "twitter", URL: "https://twitter.com/test"}

		mockSvc.On("UpdateSocialLinkService", mock.Anything, reqBody, mock.Anything, int64(5)).
			Return(nil, &repoerr.RepoError{Field: "id", Err: repoerr.ErrFrobidden}).Once()

		r := setupRouter(mockSvc, h.UpdateLinkHandler, http.MethodPut, "/profile/social/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/social/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.SocialLinkDTO{Type: "twitter", URL: "https://twitter.com/test"}

		mockSvc.On("UpdateSocialLinkService", mock.Anything, reqBody, mock.Anything, int64(5)).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.UpdateLinkHandler, http.MethodPut, "/profile/social/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/social/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== DeleteSocialLinkHandler Tests ====================

func TestDeleteSocialLinkHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteLinkService", mock.Anything, int64(5), mock.Anything).
			Return(nil).Once()

		r := setupRouter(mockSvc, h.DeleteSocialLinkHandler, http.MethodDelete, "/profile/social/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/social/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.DeleteSocialLinkHandler, http.MethodDelete, "/profile/social/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/social/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.DeleteSocialLinkHandler, http.MethodDelete, "/profile/social/:id")

		req, _ := http.NewRequest(http.MethodDelete, "/profile/social/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteLinkService", mock.Anything, int64(5), mock.Anything).
			Return(apperrors.ErrForbidden).Once()

		r := setupRouter(mockSvc, h.DeleteSocialLinkHandler, http.MethodDelete, "/profile/social/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/social/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteLinkService", mock.Anything, int64(5), mock.Anything).
			Return(errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.DeleteSocialLinkHandler, http.MethodDelete, "/profile/social/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/social/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== AddPurposesHandler Tests ====================

func TestAddPurposesHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := []dto.UserPurposeDTO{
			{Purpose: "Purpose 1"},
			{Purpose: "Purpose 2"},
		}

		expectedPurposes := []models.UserPurpose{
			{ID: 1, UserID: 1, Purpose: "Purpose 1"},
			{ID: 2, UserID: 1, Purpose: "Purpose 2"},
		}

		mockSvc.On("AddPurposesService", mock.Anything, reqBody, mock.Anything).
			Return(expectedPurposes, nil).Once()

		r := setupRouter(mockSvc, h.AddPurposesHandler, http.MethodPost, "/profile/purpose", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profile/purpose", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.AddPurposesHandler, http.MethodPost, "/profile/purpose")

		req, _ := http.NewRequest(http.MethodPost, "/profile/purpose", bytes.NewReader([]byte("[]")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.AddPurposesHandler, http.MethodPost, "/profile/purpose", 1, models.User)

		req, _ := http.NewRequest(http.MethodPost, "/profile/purpose", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := []dto.UserPurposeDTO{{Purpose: "Purpose 1"}}

		mockSvc.On("AddPurposesService", mock.Anything, reqBody, mock.Anything).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.AddPurposesHandler, http.MethodPost, "/profile/purpose", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profile/purpose", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== UpdatePurposeHandler Tests ====================

func TestUpdatePurposeHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.UserPurposeDTO{Purpose: "Updated Purpose"}

		expectedPurposes := []models.UserPurpose{
			{ID: 5, UserID: 1, Purpose: "Updated Purpose"},
		}

		mockSvc.On("UpdatePurposeService", mock.Anything, reqBody, mock.Anything, int64(5)).
			Return(expectedPurposes, nil).Once()

		r := setupRouter(mockSvc, h.UpdatePurposeHandler, http.MethodPut, "/profile/purpose/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/purpose/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.UpdatePurposeHandler, http.MethodPut, "/profile/purpose/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile/purpose/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.UpdatePurposeHandler, http.MethodPut, "/profile/purpose/:id")

		req, _ := http.NewRequest(http.MethodPut, "/profile/purpose/5", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.UpdatePurposeHandler, http.MethodPut, "/profile/purpose/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile/purpose/5", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.UserPurposeDTO{Purpose: "Updated Purpose"}

		mockSvc.On("UpdatePurposeService", mock.Anything, reqBody, mock.Anything, int64(5)).
			Return(nil, &repoerr.RepoError{Field: "id", Err: repoerr.ErrFrobidden}).Once()

		r := setupRouter(mockSvc, h.UpdatePurposeHandler, http.MethodPut, "/profile/purpose/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/purpose/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.UserPurposeDTO{Purpose: "Updated Purpose"}

		mockSvc.On("UpdatePurposeService", mock.Anything, reqBody, mock.Anything, int64(5)).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.UpdatePurposeHandler, http.MethodPut, "/profile/purpose/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/purpose/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== DeletePurposeHandler Tests ====================

func TestDeletePurposeHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeletePurposeService", mock.Anything, int64(5), mock.Anything).
			Return(nil).Once()

		r := setupRouter(mockSvc, h.DeletePurposeHandler, http.MethodDelete, "/profile/purpose/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/purpose/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.DeletePurposeHandler, http.MethodDelete, "/profile/purpose/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/purpose/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.DeletePurposeHandler, http.MethodDelete, "/profile/purpose/:id")

		req, _ := http.NewRequest(http.MethodDelete, "/profile/purpose/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeletePurposeService", mock.Anything, int64(5), mock.Anything).
			Return(apperrors.ErrForbidden).Once()

		r := setupRouter(mockSvc, h.DeletePurposeHandler, http.MethodDelete, "/profile/purpose/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/purpose/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeletePurposeService", mock.Anything, int64(5), mock.Anything).
			Return(errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.DeletePurposeHandler, http.MethodDelete, "/profile/purpose/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/purpose/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== AddProfessionsHandler Tests ====================

func TestAddProfessionsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := []dto.ProfessionDTO{
			{ProfessionID: 1},
			{ProfessionID: 2},
		}

		expectedProfessions := []models.UserProfession{
			{ID: 1, UserID: 1, ProfessionID: 1},
			{ID: 2, UserID: 1, ProfessionID: 2},
		}

		mockSvc.On("AddProfessionService", mock.Anything, mock.Anything, reqBody).
			Return(expectedProfessions, nil).Once()

		r := setupRouter(mockSvc, h.AddProfessionsHandler, http.MethodPost, "/profile/profession", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profile/profession", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.AddProfessionsHandler, http.MethodPost, "/profile/profession")

		req, _ := http.NewRequest(http.MethodPost, "/profile/profession", bytes.NewReader([]byte("[]")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.AddProfessionsHandler, http.MethodPost, "/profile/profession", 1, models.User)

		req, _ := http.NewRequest(http.MethodPost, "/profile/profession", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := []dto.ProfessionDTO{{ProfessionID: 1}}

		mockSvc.On("AddProfessionService", mock.Anything, mock.Anything, reqBody).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.AddProfessionsHandler, http.MethodPost, "/profile/profession", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPost, "/profile/profession", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== EditProfessionsHandler Tests ====================

func TestEditProfessionsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.ProfessionDTO{ProfessionID: 3}

		expectedProfession := &models.UserProfession{
			ID:           5,
			UserID:       1,
			ProfessionID: 3,
		}

		mockSvc.On("EditProfessionCategoryService", mock.Anything, mock.Anything, mock.MatchedBy(func(p *models.UserProfession) bool {
			return p.ID == 5 && p.ProfessionID == 3
		})).Return(expectedProfession, nil).Once()

		r := setupRouter(mockSvc, h.EditProfessionsHandler, http.MethodPut, "/profile/profession/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/profession/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.EditProfessionsHandler, http.MethodPut, "/profile/profession/:id")

		req, _ := http.NewRequest(http.MethodPut, "/profile/profession/5", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.EditProfessionsHandler, http.MethodPut, "/profile/profession/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile/profession/invalid", bytes.NewReader([]byte("{}")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid request body", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.EditProfessionsHandler, http.MethodPut, "/profile/profession/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodPut, "/profile/profession/5", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("forbidden", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.ProfessionDTO{ProfessionID: 3}

		mockSvc.On("EditProfessionCategoryService", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, &repoerr.RepoError{Field: "id", Err: repoerr.ErrFrobidden}).Once()

		r := setupRouter(mockSvc, h.EditProfessionsHandler, http.MethodPut, "/profile/profession/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/profession/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		reqBody := dto.ProfessionDTO{ProfessionID: 3}

		mockSvc.On("EditProfessionCategoryService", mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.EditProfessionsHandler, http.MethodPut, "/profile/profession/:id", 1, models.User)

		jsonInput, _ := json.Marshal(reqBody)
		req, _ := http.NewRequest(http.MethodPut, "/profile/profession/5", bytes.NewReader(jsonInput))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== DeleteProfessionsHandler Tests ====================

func TestDeleteProfessionsHandler(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfessionService", mock.Anything, int64(5), mock.Anything).
			Return(nil).Once()

		r := setupRouter(mockSvc, h.DeleteProfessionsHandler, http.MethodDelete, "/profile/profession/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/profession/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouter(mockSvc, h.DeleteProfessionsHandler, http.MethodDelete, "/profile/profession/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/profession/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("unauthorized", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		r := setupRouterNoAuth(h.DeleteProfessionsHandler, http.MethodDelete, "/profile/profession/:id")

		req, _ := http.NewRequest(http.MethodDelete, "/profile/profession/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockSvc := new(mocks.ProfileCore)
		h := profile_handler.NewProfileHandler(mockSvc)

		mockSvc.On("DeleteProfessionService", mock.Anything, int64(5), mock.Anything).
			Return(errors.New("database error")).Once()

		r := setupRouter(mockSvc, h.DeleteProfessionsHandler, http.MethodDelete, "/profile/profession/:id", 1, models.User)

		req, _ := http.NewRequest(http.MethodDelete, "/profile/profession/5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

// ==================== NewProfileHandler Test ====================

func TestNewProfileHandler(t *testing.T) {
	mockSvc := new(mocks.ProfileCore)
	h := profile_handler.NewProfileHandler(mockSvc)

	assert.NotNil(t, h)
}
