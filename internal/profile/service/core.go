package profileService

import (
	"context"
	"time"

	"github.com/sewaustav/CaseGoProfile/internal/profile/dto"
	"github.com/sewaustav/CaseGoProfile/internal/profile/models"
)

func (s *ProfileService) CreateProfile(
	ctx context.Context, 
	req dto.CreateProfileRequest, 
	userID int64, 
	role models.UserRole) (
		*models.Profile, 
		[]models.UserSocialLink, 
		[]models.UserPurpose, 
		error,
		) {

	tx, err := s.repo.Begin(ctx)
    if err != nil {
        return nil, nil, nil, err
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
		UserID: userID,
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
			UserID: userID,
		})
	}

	var purposes []models.UserPurpose

	for _, purpose := range req.Purposes {
		purposes = append(purposes, models.UserPurpose{
			Purpose: purpose.Purpose,
			UserID: userID,
		})
	}

	createdProfile, err := txRepo.CreateProfile(ctx, profile)
	if err != nil {
		return nil, nil, nil, err 
	}
	addedLinks, err := txRepo.AddSocial(ctx, socialLinks)
	if err != nil {
		return nil, nil, nil, err
	}
	createdPurposes, err := txRepo.AddPurposes(ctx, purposes)
	if err != nil {
		return nil, nil, nil, err 
	}

	if err := tx.Commit(); err != nil {
        return nil, nil, nil, err
    }

	return createdProfile, addedLinks, createdPurposes, nil
	
}
