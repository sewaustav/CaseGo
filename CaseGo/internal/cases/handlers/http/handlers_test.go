package http_handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	"github.com/sewaustav/CaseGoCore/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestGetCasesHandler_OK(t *testing.T) {
	svc := mocks.NewCaseGoService(t)
	h := NewCaseGoHttpHandler(svc)

	router := setupRouter()
	router.POST("/cases", h.GetCasesHandler)

	reqBody := dto.GetCasesDto{
		Limit:    10,
		Page:     1,
		Topic:    nil,
		Category: nil,
	}
	raw, _ := json.Marshal(reqBody)

	svc.On("GetCasesService", mock.Anything, 10, 1, mock.AnythingOfType("*dto.UserSettingsDto")).
		Return([]models.Case{{ID: 1}}, nil)

	req := httptest.NewRequest(http.MethodPost, "/cases", bytes.NewReader(raw))
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
}

func TestGetCaseByIDHandler_BadRequest(t *testing.T) {
	svc := mocks.NewCaseGoService(t)
	h := NewCaseGoHttpHandler(svc)

	router := setupRouter()
	router.GET("/cases/:caseID", h.GetCaseByIDHandler)

	req := httptest.NewRequest(http.MethodGet, "/cases/abc", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid case ID")
}

func TestGetCaseByIDHandler_OK(t *testing.T) {
	svc := mocks.NewCaseGoService(t)
	h := NewCaseGoHttpHandler(svc)

	router := setupRouter()
	router.GET("/cases/:caseID", h.GetCaseByIDHandler)

	svc.On("GetCaseByIDService", mock.Anything, int64(1)).Return(&models.Case{ID: 1}, nil)

	req := httptest.NewRequest(http.MethodGet, "/cases/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), `"id":1`)
}

func TestCompleteDialogHandler_BadRequest(t *testing.T) {
	svc := mocks.NewCaseGoService(t)
	h := NewCaseGoHttpHandler(svc)

	router := setupRouter()
	router.POST("/dialogs/:dialogID/complete", h.CompleteDialogHandler)

	req := httptest.NewRequest(http.MethodPost, "/dialogs/abc/complete", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Contains(t, w.Body.String(), "unauthorized")
}
