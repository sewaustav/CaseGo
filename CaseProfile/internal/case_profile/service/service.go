package service

import (
	"context"

	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/repository"
)

type Service interface {
	HandleResults(ctx context.Context, user models.UserIdentity) error
	GetProfile(ctx context.Context, user models.UserIdentity) (*models.CaseProfile, error)
	GetHistory(ctx context.Context, user models.UserIdentity) ([]*models.CaseProfileHistory, error)

	// admin only
	GetProfileByUserID(ctx context.Context, userID int64, user models.UserIdentity) (*models.CaseProfile, error)
	GetProfileByID(ctx context.Context, id int64, user models.UserIdentity) (*models.CaseProfile, error)
	GetUserHistory(ctx context.Context, userID int64, user models.UserIdentity) ([]*models.CaseProfileHistory, error)
	DeleteResultByID(ctx context.Context, id int64, user models.UserIdentity) error
}

type CaseResultService struct {
	repo repository.CaseResultRepo
}

func (c CaseResultService) HandleResults(ctx context.Context, user models.UserIdentity) error {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetProfile(ctx context.Context, user models.UserIdentity) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetHistory(ctx context.Context, user models.UserIdentity) ([]*models.CaseProfileHistory, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetProfileByUserID(ctx context.Context, userID int64, user models.UserIdentity) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetProfileByID(ctx context.Context, id int64, user models.UserIdentity) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) GetUserHistory(ctx context.Context, userID int64, user models.UserIdentity) ([]*models.CaseProfileHistory, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseResultService) DeleteResultByID(ctx context.Context, id int64, user models.UserIdentity) error {
	//TODO implement me
	panic("implement me")
}

func NewCaseResultService(repo repository.CaseResultRepo) CaseResultService {
	return CaseResultService{repo: repo}
}
