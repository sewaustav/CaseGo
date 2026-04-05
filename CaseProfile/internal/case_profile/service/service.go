package service

import (
	"context"
	"errors"
	"time"

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
	}

	if err := s.repo.AddResult(ctx, res); err != nil {
		return err
	}

	profile, err := s.repo.GetProfileByUserID(ctx, user.UserID)
	if err != nil {
		return err
	}

	if profile == nil {
		if err = s.repo.UpdateProfile(ctx, &models.CaseProfile{
			UserID:               user.UserID,
			TotalCases:           1,
			Assertiveness:        result.Assertiveness,
			Empathy:              result.Empathy,
			ClarityCommunication: result.ClarityCommunication,
			Resistance:           result.Resistance,
			Eloquence:            result.Eloquence,
			Initiative:           result.Initiative,
		}); err != nil {
			return err
		}
		return nil
	}

	newProfile, err := s.updateProfile(profile, res)

	if err = s.repo.StoreProfile(ctx, newProfile); err != nil {
		return err
	}

	if err := s.repo.UpdateProfile(ctx, newProfile); err != nil {
		return err
	}

	return nil
}

func (s CaseResultService) GetProfileService(ctx context.Context, user models.UserIdentity) (*models.CaseProfile, error) {
	return s.repo.GetProfileByUserID(ctx, user.UserID)
}

func (s CaseResultService) GetHistoryService(ctx context.Context, user models.UserIdentity, from time.Time) ([]*models.CaseProfileHistory, error) {
	return s.repo.GetHistoryBy(ctx, user.UserID, from)
}

func (s CaseResultService) GetProfileByUserIDService(ctx context.Context, userID int64, user models.UserIdentity) (*models.CaseProfile, error) {
	if user.Role != models.Admin {
		return nil, errors.New("user does not have permission to access profile by user ID")
	}
	return s.repo.GetProfileByUserID(ctx, userID)
}

func (s CaseResultService) GetProfileByIDService(ctx context.Context, id int64, user models.UserIdentity) (*models.CaseProfile, error) {
	if user.Role != models.Admin {
		return nil, errors.New("user does not have permission to access profile by ID")
	}
	return s.repo.GetProfileByID(ctx, id)
}

func (s CaseResultService) GetUserHistoryService(ctx context.Context, userID int64, user models.UserIdentity) ([]*models.CaseProfileHistory, error) {
	if user.Role != models.Admin {
		return nil, errors.New("user does not have permission to access history")
	}
	return s.repo.GetHistoryBy(ctx, userID, time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC))
}

func (s CaseResultService) DeleteResultByIDService(ctx context.Context, id int64, user models.UserIdentity) error {
	if user.Role != models.Admin {
		return errors.New("user does not have permission to delete result")
	}
	return s.repo.DeleteResultByID(ctx, id)
}

func NewCaseResultService(repo repository.CaseResultRepo) CaseResultService {
	return CaseResultService{
		repo:        repo,
		coefficient: 0.2,
	}
}
