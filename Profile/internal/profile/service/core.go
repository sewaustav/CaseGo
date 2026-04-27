package profileService

import (
	"context"
	"fmt"
	"time"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	profilerepo "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/profile_repo"
	apperr "github.com/YoungFlores/Case_Go/Profile/pkg/errors"
)

func (s *ProfileService) CreateProfileService(
	ctx context.Context,
	req dto.CreateProfileRequest,
	usr models.UserIdentity,
) (*models.UserProfile, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", apperr.ErrInternal)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	txRepo := s.repo.WithTx(tx)

	var sexPtr *models.UserSex
	if req.Info.Sex != nil {
		sexValue := models.UserSex(*req.Info.Sex)
		sexPtr = &sexValue
	}

	now := time.Now()
	profile := &models.Profile{
		UserID:      usr.UserID,
		Avatar:      req.Info.Avatar,
		IsActive:    true,
		Description: req.Info.Description,
		Username:    req.Info.Username,
		Name:        req.Info.Name,
		Surname:     req.Info.Surname,
		Patronymic:  req.Info.Patronymic,
		City:        req.Info.City,
		Age:         req.Info.Age,
		Sex:         sexPtr,
		Profession:  req.Info.Profession,
		CaseCount:   0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	var socialLinks []models.UserSocialLink
	for _, link := range req.SocialLinks {
		socialLinks = append(socialLinks, models.UserSocialLink{
			Type:   link.Type,
			URL:    link.URL,
			UserID: usr.UserID,
		})
	}

	var purposes []models.UserPurpose
	for _, purpose := range req.Purposes {
		purposes = append(purposes, models.UserPurpose{
			Purpose: purpose.Purpose,
			UserID:  usr.UserID,
		})
	}

	createdProfile, err := txRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	addedLinks, err := txRepo.AddSocial(ctx, socialLinks)
	if err != nil {
		return nil, fmt.Errorf("failed to add social links: %w", apperr.ErrInternal)
	}

	createdPurposes, err := txRepo.AddPurposes(ctx, purposes)
	if err != nil {
		return nil, fmt.Errorf("failed to add purposes: %w", apperr.ErrInternal)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", apperr.ErrInternal)
	}

	return &models.UserProfile{
		UsrProfile:  *createdProfile,
		UsrPurposes: createdPurposes,
		UsrSocials:  addedLinks,
	}, nil
}

func (s *ProfileService) UpdateProfileService(ctx context.Context, usr models.UserIdentity, req dto.ProfileInfoDTO) (*models.Profile, error) {
	now := time.Now()
	profile := &models.Profile{
		UserID:      usr.UserID,
		Avatar:      req.Avatar,
		IsActive:    true,
		Description: req.Description,
		Username:    req.Username,
		Name:        req.Name,
		Surname:     req.Surname,
		Patronymic:  req.Patronymic,
		City:        req.City,
		Age:         req.Age,
		Profession:  req.Profession,
		UpdatedAt:   now,
	}

	if req.Sex != nil {
		sexValue := models.UserSex(*req.Sex)
		profile.Sex = &sexValue
	}

	updatedProfile, err := s.repo.UpdateProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return updatedProfile, nil
}

func (s *ProfileService) UpdatePurposeService(ctx context.Context, req dto.UserPurposeDTO, usr models.UserIdentity, id int64) ([]models.UserPurpose, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func(tx profilerepo.Tx) {
		err := tx.Rollback()
		if err != nil {

		}
	}(tx)

	txRepo := s.repo.WithTx(tx)

	purpose := &models.UserPurpose{
		ID:      id,
		UserID:  usr.UserID,
		Purpose: req.Purpose,
	}

	purposes, err := txRepo.EditPurpose(ctx, purpose)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return purposes, nil
}

func (s *ProfileService) UpdateSocialLinkService(ctx context.Context, req dto.SocialLinkDTO, usr models.UserIdentity, id int64) ([]models.UserSocialLink, error) {
	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func(tx profilerepo.Tx) {
		err := tx.Rollback()
		if err != nil {

		}
	}(tx)

	txRepo := s.repo.WithTx(tx)

	link := &models.UserSocialLink{
		ID:     id,
		UserID: usr.UserID,
		Type:   req.Type,
		URL:    req.URL,
	}

	links, err := txRepo.EditSocial(ctx, link)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return links, nil

}

func (s *ProfileService) PatchProfileService(ctx context.Context, usr models.UserIdentity, req dto.UpdateProfilePartialDTO) (*models.Profile, error) {
	return s.repo.PatchProfile(ctx, usr.UserID, req)
}

func (s *ProfileService) GetUserProfileService(ctx context.Context, usr models.UserIdentity) (*models.UserProfile, error) {
	profile, err := s.repo.GetUserProfile(ctx, usr.UserID)
	if err != nil {
		return nil, err
	}

	if profile.UserID != usr.UserID {
		return nil, err
	}

	if !profile.IsActive {
		return nil, apperr.ErrIsNotActive
	}

	purposes, err := s.repo.GetUserPurposes(ctx, usr.UserID)
	if err != nil {
		return nil, err
	}

	links, err := s.repo.GetUserSocials(ctx, usr.UserID)
	if err != nil {
		return nil, err
	}

	return &models.UserProfile{
		UsrProfile:  *profile,
		UsrPurposes: purposes,
		UsrSocials:  links,
	}, nil
}

func (s *ProfileService) GetUserProfileByIDService(ctx context.Context, usr models.UserIdentity, id int64) (*models.UserProfile, error) {
	profile, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if usr.Role != models.Admin && profile.UserID != usr.UserID {
		return nil, apperr.ErrForbidden
	}

	purposes, err := s.repo.GetUserPurposes(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}

	links, err := s.repo.GetUserSocials(ctx, profile.UserID)
	if err != nil {
		return nil, err
	}

	return &models.UserProfile{
		UsrProfile:  *profile,
		UsrPurposes: purposes,
		UsrSocials:  links,
	}, nil
}

func (s *ProfileService) GetAllUsersService(ctx context.Context, usr models.UserIdentity, limit, page int) ([]models.Profile, error) {
	if usr.Role != models.Admin {
		return nil, fmt.Errorf("user is not admin")
	}

	offset := (page - 1) * limit

	users, err := s.repo.GetAllUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return users, nil

}

func (s *ProfileService) DeletePurposeService(ctx context.Context, id int64, usr models.UserIdentity) error {
	ownerID, err := s.repo.GetUserByProfileID(ctx, id, usr.UserID)
	if err != nil {
		return err
	}

	if ownerID != usr.UserID {
		return apperr.ErrForbidden
	}

	return s.repo.DeletePurpose(ctx, id)
}

func (s *ProfileService) DeleteLinkService(ctx context.Context, id int64, usr models.UserIdentity) error {
	ownerID, err := s.repo.GetUserByProfileID(ctx, id, usr.UserID)
	if err != nil {
		return err
	}

	if ownerID != usr.UserID {
		return apperr.ErrForbidden
	}

	return s.repo.DeleteSocial(ctx, id)
}

// Вспомогательные методы добавления
func (s *ProfileService) AddPurposesService(ctx context.Context, req []dto.UserPurposeDTO, usr models.UserIdentity) ([]models.UserPurpose, error) {
	var purposes []models.UserPurpose
	for _, p := range req {
		purposes = append(purposes, models.UserPurpose{Purpose: p.Purpose, UserID: usr.UserID})
	}
	return s.repo.AddPurposes(ctx, purposes)
}

func (s *ProfileService) AddSocialLinksService(ctx context.Context, req []dto.SocialLinkDTO, usr models.UserIdentity) ([]models.UserSocialLink, error) {
	var links []models.UserSocialLink
	for _, l := range req {
		links = append(links, models.UserSocialLink{UserID: usr.UserID, Type: l.Type, URL: l.URL})
	}
	return s.repo.AddSocial(ctx, links)
}

// Методы удаления профиля
func (s *ProfileService) DeleteProfileService(ctx context.Context, usr models.UserIdentity) error {
	return s.repo.DeleteProfile(ctx, usr.UserID)
}

func (s *ProfileService) DeleteProfileWithoutRecoveryService(ctx context.Context, usr models.UserIdentity, userID int64) error {
	if usr.Role != models.Admin {
		return apperr.ErrForbidden
	}
	return s.repo.DeleteProfileWithoutRecovery(ctx, userID)
}
