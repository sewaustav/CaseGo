package repository

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type CaseRepo interface {
	CreateCase(ctx context.Context, caseCopy *models.Case) (*models.Case, error)

	GetCaseByID(ctx context.Context, caseID int64) (*models.Case, error)
	GetCases(ctx context.Context, limit int, offset int) ([]models.Case, error)
	GetCasesByCategory(ctx context.Context, categoryID int32, limit int, offset int) ([]models.Case, error)
	GetCasesByTopic(ctx context.Context, topicID string, limit int, offset int) ([]models.Case, error)

	PatchCase(ctx context.Context, caseCopy *models.Case) (*models.Case, error)

	DeleteCase(ctx context.Context, caseID int64) error
}
