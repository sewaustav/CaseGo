package profilerepo

import (
	"context"
	"database/sql"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
)

type ProfileRepo interface {
	Begin(ctx context.Context) (Tx, error)
	WithTx(tx Tx) ProfileRepo

	// crud methods
	CreateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error)
	AddSocial(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error)
	AddPurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error)

	GetProfileByID(ctx context.Context, id int64) (*models.Profile, error)
	GetUserByProfileID(ctx context.Context, id, userID int64) (int64, error)
	GetUserProfile(ctx context.Context, userID int64) (*models.Profile, error)
	GetUserSocials(ctx context.Context, userID int64) ([]models.UserSocialLink, error)
	GetUserPurposes(ctx context.Context, userID int64) ([]models.UserPurpose, error)
	GetAllUsers(ctx context.Context, limit int) ([]models.Profile, error)

	UpdateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error)
	PatchProfile(ctx context.Context, userID int64, updates dto.UpdateProfilePartialDTO) (*models.Profile, error) // todo - write dto for method
	UpdateLinks(ctx context.Context, links []models.UserSocialLink) ([]models.UserSocialLink, error)
	EditSocial(ctx context.Context, link *models.UserSocialLink) ([]models.UserSocialLink, error)
	UpdatePurposes(ctx context.Context, purposes []models.UserPurpose) ([]models.UserPurpose, error)
	EditPurpose(ctx context.Context, purpose *models.UserPurpose) ([]models.UserPurpose, error)

	DeletePurpose(ctx context.Context, id int64) error
	DeleteSocial(ctx context.Context, id int64) error
	DeleteProfile(ctx context.Context, userID int64) error
	DeleteProfileWithoutRecovery(ctx context.Context, userID int64) error

	// professions methods
	AddProfession(ctx context.Context, profession *models.UserProfession) (*models.UserProfession, error)
	GetAllProfessions(ctx context.Context, userID int64) ([]models.UserProfession, error)
	GetProfileIDByProfessionID(ctx context.Context, professionID int64) (int64, error)
	EditProfession(ctx context.Context, profession *models.UserProfession) (*models.UserProfession, error)
	DeleteProfession(ctx context.Context, professionID int64) error
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

type PostgresProfileRepo struct {
	db DBTX
}

func NewPostgresProfileRepo(db *sql.DB) *PostgresProfileRepo {
	return &PostgresProfileRepo{
		db: db,
	}
}
