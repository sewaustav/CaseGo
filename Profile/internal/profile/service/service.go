package profileService

import (
	"context"

	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/repo"
	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	profilerepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
)

type ProfileCore interface {
	CreateProfileService(ctx context.Context, req dto.CreateProfileRequest, usr models.UserIdentity) (*models.UserProfile, error)
	UpdateProfileService(ctx context.Context, usr models.UserIdentity, req dto.ProfileInfoDTO) (*models.Profile, error)
	PatchProfileService(ctx context.Context, usr models.UserIdentity, req dto.UpdateProfilePartialDTO) (*models.Profile, error)
	GetUserProfileService(ctx context.Context, usr models.UserIdentity) (*models.UserProfile, error)
	GetUserProfileByIDService(ctx context.Context, usr models.UserIdentity, id int64) (*models.UserProfile, error)
	DeleteProfileService(ctx context.Context, usr models.UserIdentity) error
	DeleteProfileWithoutRecoveryService(ctx context.Context, usr models.UserIdentity, userID int64) error

	AddPurposesService(ctx context.Context, req []dto.UserPurposeDTO, usr models.UserIdentity) ([]models.UserPurpose, error)
	UpdatePurposeService(ctx context.Context, req dto.UserPurposeDTO, usr models.UserIdentity, id int64) ([]models.UserPurpose, error)
	DeletePurposeService(ctx context.Context, id int64, usr models.UserIdentity) error

	AddSocialLinksService(ctx context.Context, req []dto.SocialLinkDTO, usr models.UserIdentity) ([]models.UserSocialLink, error)
	UpdateSocialLinkService(ctx context.Context, req dto.SocialLinkDTO, usr models.UserIdentity, id int64) ([]models.UserSocialLink, error)
	DeleteLinkService(ctx context.Context, id int64, usr models.UserIdentity) error

	AddProfessionService(ctx context.Context, usr models.UserIdentity, req []dto.ProfessionDTO) ([]models.UserProfession, error)
	EditProfessionCategoryService(ctx context.Context, usr models.UserIdentity, profession *models.UserProfession) (*models.UserProfession, error)
	DeleteProfessionService(ctx context.Context, id int64, usr models.UserIdentity) error
	GetProfessionsService(ctx context.Context, usr models.UserIdentity) ([]models.UserProfession, error)
}

type ProfileService struct {
	repo    profilerepo.ProfileRepo
	catRepo repo.CategoryRepo
}

func NewProfileService(repo profilerepo.ProfileRepo, catRepo repo.CategoryRepo) *ProfileService {
	return &ProfileService{
		repo:    repo,
		catRepo: catRepo,
	}
}
