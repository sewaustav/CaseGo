package profilerepo

import (
	"context"
	"database/sql"

	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
)

type ProfileRepo interface {
	CreateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error)
	AddSocial(ctx context.Context, links []string) error
	CreatePupose(ctx context.Context, pupose string) (*models.UserPurpose, error)
	
	GetUserProfile(ctx context.Context, userId int64) (*models.Profile, error)
	GetUserSocials(ctx context.Context, userID int64) ([]string, error)
	GetUserPurposes(ctx context.Context, userID int64)
	GetAllUsers(ctx context.Context, limit int) ([]models.Profile, error)

	UpdateProfile(ctx context.Context, user *models.Profile) (*models.Profile, error)
	EditSocial(ctx context.Context, link string, id int64) ([]string, error)
	EditPurpose(ctx context.Context, purpose string, id int64) (*models.UserPurpose, error)

	DeletePupose(ctx context.Context, id int64) error
	DeleteSocial(ctx context.Context, id int64) error
	DeleteProfile(ctx context.Context, userID int64) error
	DeleteProfileWithoutRecovery(ctx context.Context, userID int64) error

}

type PostgresProfileRepo struct {
	db *sql.DB
}

func NewPostgresProfileRepo(db *sql.DB) *PostgresProfileRepo {
	return &PostgresProfileRepo{
		db: db,
	}
}