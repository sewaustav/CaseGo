package cache

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type Interactor interface {
	Push(ctx context.Context, inter *models.Interaction) error
	GetFullHistory(ctx context.Context, dialogID int64) ([]models.Interaction, error)
	DeleteLast(ctx context.Context, dialogID int64) error
	Clear(ctx context.Context, dialogID int64) error
	Close() error
}
