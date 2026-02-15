package profileService

import (
	"github.com/YoungFlores/Case_Go/Profile/internal/profession_categories/repo"
	profilerepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
)

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
