package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/sewaustav/CaseGoProfile/apperrors"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProfileService_Success(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}
	expected := &models.CaseProfile{ID: 1, UserID: 1, TotalCases: 5}

	repo.On("GetProfileByUserID", ctx, int64(1)).Return(expected, nil)

	result, err := svc.GetProfileService(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetProfileService_RepoError(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}

	repo.On("GetProfileByUserID", ctx, int64(1)).Return(nil, errors.New("db error"))

	result, err := svc.GetProfileService(ctx, user)

	assert.Error(t, err)
	assert.Nil(t, result)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusInternalServerError, appErr.Code)
}

func TestGetHistoryService_Success(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}
	from := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := []*models.CaseResult{
		{ID: 1, UserID: 1, Assertiveness: 0.5},
	}

	repo.On("GetResultsByUserID", ctx, int64(1)).Return(expected, nil)

	result, err := svc.GetHistoryService(ctx, user, from)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetProfileByUserIDService_Admin_Success(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	admin := models.UserIdentity{UserID: 99, Role: models.Admin}
	expected := &models.CaseProfile{ID: 1, UserID: 42, TotalCases: 10}

	repo.On("GetProfileByUserID", ctx, int64(42)).Return(expected, nil)

	result, err := svc.GetProfileByUserIDService(ctx, 42, admin)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetProfileByUserIDService_NonAdmin_Forbidden(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}

	result, err := svc.GetProfileByUserIDService(ctx, 42, user)

	assert.Error(t, err)
	assert.Nil(t, result)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

func TestGetProfileByIDService_Admin_Success(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	admin := models.UserIdentity{UserID: 99, Role: models.Admin}
	expected := &models.CaseProfile{ID: 10, UserID: 42}

	repo.On("GetProfileByID", ctx, int64(10)).Return(expected, nil)

	result, err := svc.GetProfileByIDService(ctx, 10, admin)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetProfileByIDService_NonAdmin_Forbidden(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}

	result, err := svc.GetProfileByIDService(ctx, 10, user)

	assert.Error(t, err)
	assert.Nil(t, result)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

func TestGetUserHistoryService_Admin_Success(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	admin := models.UserIdentity{UserID: 99, Role: models.Admin}
	expected := []*models.CaseProfileHistory{
		{ID: 1, UserID: 42},
	}

	repo.On("GetHistoryBy", ctx, int64(42), time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)).Return(expected, nil)

	result, err := svc.GetUserHistoryService(ctx, 42, admin)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetUserHistoryService_NonAdmin_Forbidden(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}

	result, err := svc.GetUserHistoryService(ctx, 42, user)

	assert.Error(t, err)
	assert.Nil(t, result)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

func TestDeleteResultByIDService_Admin_Success(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	admin := models.UserIdentity{UserID: 99, Role: models.Admin}

	repo.On("DeleteResultByID", ctx, int64(5)).Return(nil)

	err := svc.DeleteResultByIDService(ctx, 5, admin)

	assert.NoError(t, err)
}

func TestDeleteResultByIDService_NonAdmin_Forbidden(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}

	err := svc.DeleteResultByIDService(ctx, 5, user)

	assert.Error(t, err)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

func TestHandleResultsService_AddResultError(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}
	result := dto.Result{
		CaseID:        1,
		DialogID:      1,
		StepsCount:    5,
		TokensUsed:    100,
		Assertiveness: 0.8,
		Empathy:       0.7,
	}

	repo.On("AddResult", ctx, mock.AnythingOfType("*models.CaseResult")).Return(errors.New("insert error"))

	err := svc.HandleResultsService(ctx, result, user)

	assert.Error(t, err)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusInternalServerError, appErr.Code)
}

func TestHandleResultsService_NewProfile(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}
	result := dto.Result{
		CaseID:               1,
		DialogID:             10,
		StepsCount:           5,
		TokensUsed:           100,
		Assertiveness:        0.8,
		Empathy:              0.7,
		ClarityCommunication: 0.6,
		Resistance:           0.5,
		Eloquence:            0.4,
		Initiative:           0.3,
	}

	repo.On("AddResult", ctx, mock.AnythingOfType("*models.CaseResult")).Return(nil)
	// Репо теперь возвращает ErrNotFound вместо nil, nil
	repo.On("GetProfileByUserID", ctx, int64(1)).Return(nil, fmt.Errorf("profile user_id=1: %w", apperrors.ErrNotFound))
	repo.On("UpdateProfile", ctx, mock.MatchedBy(func(p *models.CaseProfile) bool {
		return p.UserID == 1 && p.TotalCases == 1 && p.Assertiveness == float32(0.8)
	})).Return(nil)

	err := svc.HandleResultsService(ctx, result, user)

	assert.NoError(t, err)
}

func TestHandleResultsService_GetProfileError(t *testing.T) {
	repo := mocks.NewCaseResultRepo(t)
	svc := NewCaseResultService(repo)

	ctx := context.Background()
	user := models.UserIdentity{UserID: 1, Role: models.User}
	result := dto.Result{CaseID: 1, DialogID: 10}

	repo.On("AddResult", ctx, mock.AnythingOfType("*models.CaseResult")).Return(nil)
	repo.On("GetProfileByUserID", ctx, int64(1)).Return(nil, errors.New("db error"))

	err := svc.HandleResultsService(ctx, result, user)

	assert.Error(t, err)

	var appErr *apperrors.AppError
	assert.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusInternalServerError, appErr.Code)
}
