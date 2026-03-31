package service

import (
	"context"
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
	repo repository.CaseResultRepo
}

func (c CaseResultService) HandleResultsService(ctx context.Context, result dto.Result, user models.UserIdentity) error {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetProfileService(ctx context.Context, user models.UserIdentity) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetHistoryService(ctx context.Context, user models.UserIdentity, from time.Time) ([]*models.CaseProfileHistory, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetProfileByUserIDService(ctx context.Context, userID int64, user models.UserIdentity) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetProfileByIDService(ctx context.Context, id int64, user models.UserIdentity) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetUserHistoryService(ctx context.Context, userID int64, user models.UserIdentity) ([]*models.CaseProfileHistory, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) DeleteResultByIDService(ctx context.Context, id int64, user models.UserIdentity) error {
	//TODO implement me
	panic("implement me")
}

func NewCaseResultService(repo repository.CaseResultRepo) CaseResultService {
	return CaseResultService{repo: repo}
}
