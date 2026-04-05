package service

import (
	"context"
	"errors"
	"time"

	"github.com/sewaustav/CaseGoProfile/apperrors"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/repository"
)

type Service interface {
	HandleResultsService(ctx context.Context, result dto.Result, user models.UserIdentity) error
	GetProfileService(ctx context.Context, user models.UserIdentity) (*models.CaseProfile, error)
	GetHistoryService(ctx context.Context, user models.UserIdentity, from time.Time) ([]*models.CaseProfileHistory, error)

	// admin only
	GetProfileByUserIDService(ctx context.Context, userID int64, user models.UserIdentity) (*models.CaseProfile, error)
	GetProfileByIDService(ctx context.Context, id int64, user models.UserIdentity) (*models.CaseProfile, error)
	GetUserHistoryService(ctx context.Context, userID int64, user models.UserIdentity) ([]*models.CaseProfileHistory, error)
	DeleteResultByIDService(ctx context.Context, id int64, user models.UserIdentity) error
}

type CaseResultService struct {
	repo        repository.CaseResultRepo
	coefficient float32
}

func NewCaseResultService(repo repository.CaseResultRepo) CaseResultService {
	return CaseResultService{
		repo:        repo,
		coefficient: 0.2,
	}
}

func (s CaseResultService) HandleResultsService(ctx context.Context, result dto.Result, user models.UserIdentity) error {
	res := &models.CaseResult{
		UserID:               user.UserID,
		CaseID:               result.CaseID,
		DialogID:             result.DialogID,
		StepsCount:           result.StepsCount,
		TokensUsed:           result.TokensUsed,
		Assertiveness:        result.Assertiveness,
		Empathy:              result.Empathy,
		ClarityCommunication: result.ClarityCommunication,
		Resistance:           result.Resistance,
		Eloquence:            result.Eloquence,
		Initiative:           result.Initiative,
	}

	if err := s.repo.AddResult(ctx, res); err != nil {
		return apperrors.NewInternal("failed to add case result", err)
	}

	profile, err := s.repo.GetProfileByUserID(ctx, user.UserID)
	if err != nil {
		// Профиль не найден — первый кейс пользователя, создаём новый
		if errors.Is(err, apperrors.ErrNotFound) {
			if err := s.repo.UpdateProfile(ctx, &models.CaseProfile{
				UserID:               user.UserID,
				TotalCases:           1,
				Assertiveness:        result.Assertiveness,
				Empathy:              result.Empathy,
				ClarityCommunication: result.ClarityCommunication,
				Resistance:           result.Resistance,
				Eloquence:            result.Eloquence,
				Initiative:           result.Initiative,
			}); err != nil {
				return apperrors.NewInternal("failed to create initial profile", err)
			}
			return nil
		}
		return apperrors.NewInternal("failed to get profile", err)
	}

	newProfile, _ := s.updateProfile(profile, res)

	if err := s.repo.StoreProfile(ctx, newProfile); err != nil {
		return apperrors.NewInternal("failed to store profile history", err)
	}

	if err := s.repo.UpdateProfile(ctx, newProfile); err != nil {
		return apperrors.NewInternal("failed to update profile", err)
	}

	return nil
}

func (s CaseResultService) GetProfileService(ctx context.Context, user models.UserIdentity) (*models.CaseProfile, error) {
	profile, err := s.repo.GetProfileByUserID(ctx, user.UserID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("case profile not found", err)
		}
		return nil, apperrors.NewInternal("failed to get profile", err)
	}
	return profile, nil
}

func (s CaseResultService) GetHistoryService(ctx context.Context, user models.UserIdentity, from time.Time) ([]*models.CaseProfileHistory, error) {
	history, err := s.repo.GetHistoryBy(ctx, user.UserID, from)
	if err != nil {
		return nil, apperrors.NewInternal("failed to get history", err)
	}
	return history, nil
}

func (s CaseResultService) GetProfileByUserIDService(ctx context.Context, userID int64, user models.UserIdentity) (*models.CaseProfile, error) {
	if user.Role != models.Admin {
		return nil, apperrors.NewForbidden("only admin can access profile by user ID", nil)
	}

	profile, err := s.repo.GetProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("case profile not found", err)
		}
		return nil, apperrors.NewInternal("failed to get profile by user ID", err)
	}
	return profile, nil
}

func (s CaseResultService) GetProfileByIDService(ctx context.Context, id int64, user models.UserIdentity) (*models.CaseProfile, error) {
	if user.Role != models.Admin {
		return nil, apperrors.NewForbidden("only admin can access profile by ID", nil)
	}

	profile, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("case profile not found", err)
		}
		return nil, apperrors.NewInternal("failed to get profile by ID", err)
	}
	return profile, nil
}

func (s CaseResultService) GetUserHistoryService(ctx context.Context, userID int64, user models.UserIdentity) ([]*models.CaseProfileHistory, error) {
	if user.Role != models.Admin {
		return nil, apperrors.NewForbidden("only admin can access user history", nil)
	}

	history, err := s.repo.GetHistoryBy(ctx, userID, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		return nil, apperrors.NewInternal("failed to get user history", err)
	}
	return history, nil
}

func (s CaseResultService) DeleteResultByIDService(ctx context.Context, id int64, user models.UserIdentity) error {
	if user.Role != models.Admin {
		return apperrors.NewForbidden("only admin can delete results", nil)
	}

	err := s.repo.DeleteResultByID(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NewNotFound("result not found", err)
		}
		return apperrors.NewInternal("failed to delete result", err)
	}
	return nil
}
