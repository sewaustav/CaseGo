package profilerepo

import (
	"context"
	"database/sql"

	"github.com/sewaustav/CaseGoProfile/internal/profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
)

type ProfileRepo interface {
	Begin(ctx context.Context) (*sql.Tx, error)
    WithTx(tx *sql.Tx) ProfileRepo

	// crud methods 
	CreateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error)
	AddSocial(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error)
	AddPurposes(ctx context.Context, puposes []models.UserPurpose) ([]models.UserPurpose, error)
	
	GetProfileByID(ctx context.Context, id int64) (*models.Profile, error)
	GetUserProfile(ctx context.Context, userID int64) (*models.Profile, error)
	GetUserSocials(ctx context.Context, userID int64) ([]models.UserSocialLink, error)
	GetUserPurposes(ctx context.Context, userID int64) ([]models.UserPurpose, error)
	GetAllUsers(ctx context.Context, limit int) ([]models.Profile, error)

	UpdateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error)
	PathcProfile(ctx context.Context, userID int64, updates dto.UpdateProfilePartialDTO) (*models.Profile, error) // todo - write dto for method
	UpdateLinks(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error) 
	EditSocial(ctx context.Context, link *models.UserSocialLink) ([]models.UserSocialLink, error)
	UpdatePurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error)
	EditPurpose(ctx context.Context, purpose *models.UserPurpose) ([]models.UserPurpose, error)

	DeletePupose(ctx context.Context, id int64) error
	DeleteSocial(ctx context.Context, id int64) error
	DeleteProfile(ctx context.Context, userID int64) error
	DeleteProfileWithoutRecovery(ctx context.Context, userID int64) error

}

type DBTX interface {
    ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
    PrepareContext(context.Context, string) (*sql.Stmt, error)
    QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
    QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type PostgresProfileRepo struct {
	db DBTX
}

func NewPostgresProfileRepo(db *sql.DB) *PostgresProfileRepo {
	return &PostgresProfileRepo{
		db: db,
	}
}