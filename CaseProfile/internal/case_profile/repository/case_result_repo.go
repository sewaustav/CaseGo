package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
)

type CaseResultRepo interface {
	AddResult(ctx context.Context, result *models.CaseProfile) error
	GetProfileByUserID(ctx context.Context, userID int64) (*models.CaseProfile, error)
	GetProfileByID(ctx context.Context, id int64) (*models.CaseProfile, error)
	GetHistoryBy(ctx context.Context, userID int64, from time.Time) ([]*models.CaseProfileHistory, error)
	DeleteResultByID(ctx context.Context, id int64) error
}

type PostgresCaseResultRepo struct {
	db *sql.DB
}

func NewPostgresCaseResultRepo(db *sql.DB) CaseResultRepo {
	return &PostgresCaseResultRepo{db: db}
}

func (p PostgresCaseResultRepo) AddResult(ctx context.Context, result *models.CaseProfile) error {
	//TODO implement me
	panic("implement me")
}

func (p PostgresCaseResultRepo) GetProfileByUserID(ctx context.Context, userID int64) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresCaseResultRepo) GetProfileByID(ctx context.Context, id int64) (*models.CaseProfile, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresCaseResultRepo) GetHistoryBy(ctx context.Context, userID int64, from time.Time) ([]*models.CaseProfileHistory, error) {
	//TODO implement me
	panic("implement me")
}

func (p PostgresCaseResultRepo) DeleteResultByID(ctx context.Context, id int64) error {
	//TODO implement me
	panic("implement me")
}
