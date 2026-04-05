// CaseProfile/internal/case_profile/handlers/http/handler_test.go
package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// helper to set user identity in gin context
func setUserContext(c *gin.Context, userID int64, role int) {
	c.Set(UserIDKey, userID)
	c.Set(RoleKey, role)
}

// ======================== GetProfileHandler ========================

func TestGetProfileHandler_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	expected := &models.CaseProfile{ID: 1, UserID: 42, TotalCases: 5, Assertiveness: 0.8}
	svc.On("GetProfileService", mock.Anything, models.UserIdentity{UserID: 42, Role: models.User}).
		Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/profile", nil)
	setUserContext(c, 42, int(models.User))

	handler.GetProfileHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.CaseProfile
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(42), resp.UserID)
}

func TestGetProfileHandler_Unauthorized(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/profile", nil)
	// не устанавливаем user context

	handler.GetProfileHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetProfileHandler_ServiceError(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	svc.On("GetProfileService", mock.Anything, models.UserIdentity{UserID: 1, Role: models.User}).
		Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/profile", nil)
	setUserContext(c, 1, int(models.User))

	handler.GetProfileHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ======================== GetHistoryHandler ========================

func TestGetHistoryHandler_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	from := time.Date(2026, 3, 1, 0, 0, 0, 0, time.UTC)
	expected := []*models.CaseProfileHistory{
		{ID: 1, UserID: 10, Assertiveness: 0.5, Date: from},
	}

	svc.On("GetHistoryService", mock.Anything, models.UserIdentity{UserID: 10, Role: models.User}, from).
		Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/history?from=2026-03-01", nil)
	setUserContext(c, 10, int(models.User))

	handler.GetHistoryHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetHistoryHandler_DefaultFromDate(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	defaultFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	svc.On("GetHistoryService", mock.Anything, models.UserIdentity{UserID: 10, Role: models.User}, defaultFrom).
		Return([]*models.CaseProfileHistory{}, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/history", nil)
	setUserContext(c, 10, int(models.User))

	handler.GetHistoryHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetHistoryHandler_InvalidFromDate(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/history?from=invalid-date", nil)
	setUserContext(c, 10, int(models.User))

	handler.GetHistoryHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetHistoryHandler_Unauthorized(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/history", nil)

	handler.GetHistoryHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ======================== GetUserProfileHandler ========================

func TestGetUserProfileHandler_ByUserID_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	expected := &models.CaseProfile{ID: 1, UserID: 42}

	svc.On("GetProfileByUserIDService", mock.Anything, int64(42), models.UserIdentity{UserID: 99, Role: models.Admin}).
		Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/profile", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "42"}}
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserProfileHandler_ByID_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	expected := &models.CaseProfile{ID: 10, UserID: 42}

	svc.On("GetProfileByIDService", mock.Anything, int64(10), models.UserIdentity{UserID: 99, Role: models.Admin}).
		Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/profile", nil)
	c.Params = gin.Params{{Key: "id", Value: "10"}}
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserProfileHandler_NoParams(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/profile", nil)
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserProfileHandler_InvalidUserID(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/profile", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "abc"}}
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserProfileHandler_Unauthorized(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/profile", nil)

	handler.GetUserProfileHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ======================== GetUserProfileHistoryHandler ========================

func TestGetUserProfileHistoryHandler_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	expected := []*models.CaseProfileHistory{
		{ID: 1, UserID: 42},
	}

	svc.On("GetUserHistoryService", mock.Anything, int64(42), models.UserIdentity{UserID: 99, Role: models.Admin}).
		Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/history/42", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "42"}}
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHistoryHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetUserProfileHistoryHandler_NoUserID(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/history/", nil)
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHistoryHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserProfileHistoryHandler_InvalidUserID(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/history/abc", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "abc"}}
	setUserContext(c, 99, int(models.Admin))

	handler.GetUserProfileHistoryHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetUserProfileHistoryHandler_Unauthorized(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/admin/history/42", nil)

	handler.GetUserProfileHistoryHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// ======================== DeleteResultByIDHandler ========================

func TestDeleteResultByIDHandler_Success(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	svc.On("DeleteResultByIDService", mock.Anything, int64(5), models.UserIdentity{UserID: 99, Role: models.Admin}).
		Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/admin/result/5", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "5"}}
	setUserContext(c, 99, int(models.Admin))

	handler.DeleteResultByIDHandler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteResultByIDHandler_InvalidID(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/admin/result/abc", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "abc"}}
	setUserContext(c, 99, int(models.Admin))

	handler.DeleteResultByIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteResultByIDHandler_Unauthorized(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/admin/result/5", nil)

	handler.DeleteResultByIDHandler(c)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestDeleteResultByIDHandler_ServiceError(t *testing.T) {
	svc := mocks.NewService(t)
	handler := NewHttpHandler(svc)

	svc.On("DeleteResultByIDService", mock.Anything, int64(5), models.UserIdentity{UserID: 99, Role: models.Admin}).
		Return(errors.New("forbidden"))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, "/admin/result/5", nil)
	c.Params = gin.Params{{Key: "user_id", Value: "5"}}
	setUserContext(c, 99, int(models.Admin))

	handler.DeleteResultByIDHandler(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
