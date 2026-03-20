package repository

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type DialogRepo interface {
	StartDialog(ctx context.Context, dialog *models.Dialog) (*models.Dialog, error)

	GetDialogByID(ctx context.Context, dialogID int64) (*models.Dialog, error)
	GetUserDialogs(ctx context.Context, userID int64, limit, offset int) ([]models.Dialog, error)
	GetDialogsByCaseID(ctx context.Context, caseID int64, limit, offset int) ([]models.Dialog, error)

	CompleteDialog(ctx context.Context, dialogID int64) (*models.Dialog, error)
}

type Interaction interface {
	PutInteraction(ctx context.Context, interaction *models.Interaction) (*models.Interaction, error)

	GetInteractionByID(ctx context.Context, interactionID int64) (*models.Interaction, error)
	GetInteractionsByDialogID(ctx context.Context, dialogID int64) ([]models.Interaction, error)
}
