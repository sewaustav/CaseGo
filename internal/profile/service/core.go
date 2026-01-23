package profileService

import (
	"context"
	"fmt"

	"time"

	"github.com/sewaustav/CaseGoProfile/internal/profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
)

func (s *ProfileService) CreateProfileService(
	ctx context.Context, 
	req dto.CreateProfileRequest, 
	usr models.UserIdentity) (
		*models.UserProfile, 
		error,
		) {

	tx, err := s.repo.Begin(ctx)
    if err != nil {
        return nil, err
    }

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)
	
	var sexPtr *models.UserSex
	if req.Info.Sex != nil {
		sexValue := models.UserSex(*req.Info.Sex)
		sexPtr = &sexValue
	}

	now := time.Now()
	
	profile := &models.Profile{
		UserID: usr.UserID,
		Avatar: req.Info.Avatar,
		IsActive: true,
		Description: req.Info.Description,
		Username: req.Info.Username,
		Name: req.Info.Name,
		Surname: req.Info.Surname,
		Email: req.Info.Email,
		CaseCount: 0,
		CreatedAt: now,
		UpdatedAt: now,
		// optional - nil possible
		Patronymic: req.Info.Patronymic,
		PhoneNumber: req.Info.PhoneNumber,
		Sex: sexPtr,
		Profession: req.Info.Profession,
	}

	var socialLinks []models.UserSocialLink

	for _, link := range req.SocialLinks {
		socialLinks = append(socialLinks, models.UserSocialLink{
			Type: link.Type, 
			URL: link.URL, 
			UserID: usr.UserID,
		})
	}

	var purposes []models.UserPurpose

	for _, purpose := range req.Purposes {
		purposes = append(purposes, models.UserPurpose{
			Purpose: purpose.Purpose,
			UserID: usr.UserID,
		})
	}

	createdProfile, err := txRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, err 
	}
	addedLinks, err := txRepo.AddSocial(ctx, socialLinks)
	if err != nil {
		return nil, err
	}
	createdPurposes, err := txRepo.AddPurposes(ctx, purposes)
	if err != nil {
		return nil, err 
	}

	if err := tx.Commit(); err != nil {
        return nil, err
    }

	return &models.UserProfile{
		UsrProfile: *createdProfile,
		UsrPurposes: createdPurposes,
		UsrSocials: addedLinks,
	}, nil
	
}

func (s *ProfileService) AddPurposesService(ctx context.Context, req []dto.UserPurposeDTO, usr models.UserIdentity) ([]models.UserPurpose, error) {
	
	var purposes []models.UserPurpose

	for _, purpose := range req {
		purposes = append(purposes, models.UserPurpose{
			Purpose: purpose.Purpose,
			UserID: usr.UserID,
		})
	}

	createdPurposes, err := s.repo.AddPurposes(ctx, purposes) 
	if err != nil {
		return nil, err 
	}

	return createdPurposes, nil 
}

func (s *ProfileService) AddSocialLinksService(ctx context.Context, req []dto.SocialLinkDTO, usr models.UserIdentity) ([]models.UserSocialLink, error) {
	var links []models.UserSocialLink 

	for _, link := range req {
		links = append(links, models.UserSocialLink{
			UserID: usr.UserID,
			Type: link.Type,
			URL: link.URL,
		})
	}

	addedLinks, err := s.repo.AddSocial(ctx, links)
	if err != nil {
		return nil, err 
	}

	return addedLinks, nil 
} 

// method for put
func(s *ProfileService) UpdateProfileService(ctx context.Context, usr models.UserIdentity, req dto.ProfileInfoDTO) (*models.Profile, error) {
	var sexPtr *models.UserSex
	if req.Sex != nil {
		sexValue := models.UserSex(*req.Sex)
		sexPtr = &sexValue
	}

	// todo - check is username unique

	now := time.Now()
	
	profile := &models.Profile{
		UserID: usr.UserID,
		Avatar: req.Avatar,
		IsActive: true,
		Description: req.Description,
		Username: req.Username,
		Name: req.Name,
		Surname: req.Surname,
		Email: req.Email,
		CaseCount: 0,
		UpdatedAt: now,
		// optional - nil possible
		Patronymic: req.Patronymic,
		PhoneNumber: req.PhoneNumber,
		Sex: sexPtr,
		Profession: req.Profession,
	}

	updatedProfile, err := s.repo.UpdateProfile(ctx, profile) 
	if err != nil {
		return nil, err 
	}

	return updatedProfile, nil
}


// partial - update. Note - in future special method for email/phone update
func (s *ProfileService) PatchProfileService(
	ctx context.Context,
	usr models.UserIdentity,
	req dto.UpdateProfilePartialDTO,
) (*models.Profile, error) {

	// todo - check is username unique


	profile, err := s.repo.PathcProfile(ctx, usr.UserID, req)

	if err != nil {
		return nil, err 
	}

	return profile, nil 
}

func (s *ProfileService) UpdateSocialLinkService(ctx context.Context, req dto.SocialLinkDTO, usr models.UserIdentity, id int64) ([]models.UserSocialLink, error) {
	link := &models.UserSocialLink{
		ID: id,
		UserID: usr.UserID,
		Type: req.Type,
		URL: req.URL,
	}
	
	links, err := s.repo.EditSocial(ctx, link)

	if err != nil {
		return nil, err 
	}

	return links, nil

}

func (s *ProfileService) UpdatePurposeService(ctx context.Context, req dto.UserPurposeDTO, usr models.UserIdentity, id int64) ([]models.UserPurpose, error) {
	purpose := &models.UserPurpose{
		ID: id,
		UserID: usr.UserID,
		Purpose: req.Purpose,
	}

	purposes, err := s.repo.EditPurpose(ctx, purpose)
	if err != nil {
		return nil, err 
	}

	return purposes, nil 
}


// get methods
func (s *ProfileService) GetUserProfileService(ctx context.Context, usr models.UserIdentity) (*models.UserProfile, error) {
	profile, err := s.repo.GetUserProfile(ctx, usr.UserID)
	if err != nil {
		return nil, err 
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
		UsrProfile: *profile,
		UsrPurposes: purposes,
		UsrSocials: links,
	}, nil
}

func (s *ProfileService) GetUserProfileByIDService(ctx context.Context, usr models.UserIdentity, id int64) (*models.UserProfile, error) {
	
	profile, err := s.repo.GetProfileByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if usr.Role != models.Admin || profile.UserID != usr.UserID {
		return nil, fmt.Errorf("forbidden")
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
		UsrProfile: *profile,
		UsrPurposes: purposes,
		UsrSocials: links,
	}, nil 

}

func (s *ProfileService) DeleteProfileService(ctx context.Context, usr models.UserIdentity) error {
	
	if err := s.repo.DeleteProfile(ctx, usr.UserID); err != nil {
		return err 
	}

	return nil
}

func (s *ProfileService) DeleteProfileWithoutRecoveryService(ctx context.Context, usr models.UserIdentity) error {
	if usr.Role != models.Admin {
		return fmt.Errorf("forbidden")
	}

	if err := s.repo.DeleteProfileWithoutRecovery(ctx, usr.UserID); err != nil {
		return err 
	}

	return nil 
}

func (s *ProfileService) DeletePuposeService(ctx context.Context, id int64, usr models.UserIdentity) error {
	userID, err := s.repo.GetUserByProfileID(ctx, id, usr.UserID)
	if err != nil {
		return err 
	}

	if userID != usr.UserID {
		return fmt.Errorf("forbidden")
	}

	if err := s.repo.DeletePupose(ctx, id); err != nil {
		return err 
	}

	return nil 
}

func (s *ProfileService) DeleteLinkService(ctx context.Context, id int64, usr models.UserIdentity) error {
	userID, err := s.repo.GetUserByProfileID(ctx, id, usr.UserID)
	if err != nil {
		return err 
	}

	if userID != usr.UserID {
		return fmt.Errorf("forbidden")
	}

	if err := s.repo.DeleteSocial(ctx, id); err != nil {
		return err 
	}

	return nil 
}
