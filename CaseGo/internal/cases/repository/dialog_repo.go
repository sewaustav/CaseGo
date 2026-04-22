package repository

import (
	"context"
	"database/sql"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type DialogRepo interface {
	StartDialog(ctx context.Context, dialog *models.Dialog) (*models.Dialog, error)

	GetDialogByID(ctx context.Context, dialogID int64) (*models.Dialog, error)
	GetUserDialogs(ctx context.Context, userID int64, limit, offset int) ([]models.Dialog, error)
	GetDialogsByCaseID(ctx context.Context, caseID int64, limit, offset int) ([]models.Dialog, error)
	CountDialogs(ctx context.Context) (int, error)
}

type Interaction interface {
	Begin(ctx context.Context) (Tx, error)
	WithTx(tx Tx) Interaction

	PutInteraction(ctx context.Context, interaction *models.Interaction) error

	GetInteractionByID(ctx context.Context, interactionID int64) (*models.Interaction, error)
	GetInteractionsByDialogID(ctx context.Context, dialogID int64) ([]models.Interaction, error)
}

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Tx interface {
	DBTX
	Commit() error
	Rollback() error
}
