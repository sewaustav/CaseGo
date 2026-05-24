package profileService

import (
	"context"
	"errors"
	"fmt"
	"math"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
)

func (s *ProfileService) UpdateLevelService(ctx context.Context, usr models.UserIdentity, req *dto.LevelDto) (*models.UserLevel, error) {
	currentProfile, err := s.levelRepo.GetUserLevel(ctx, usr.UserID)
	if err != nil && !errors.Is(err, repoerr.ErrNotFound) {
		return nil, fmt.Errorf("failed to get user level: %w", err)
	}

	if currentProfile == nil {
		updatedProfile, err := s.levelRepo.CreateLevel(ctx, &models.UserLevel{
			UserID:     usr.UserID,
			Xp:         req.Xp,
			Level:      int(math.Floor(math.Sqrt(float64(req.Xp/100))) + 1),
			Streak:     0,
			LastActive: req.Date,
		})

		if err != nil {
			return nil, fmt.Errorf("failed to create user level: %w", err)
		}

		return updatedProfile, nil
	}

	level := math.Floor(math.Sqrt(float64((req.Xp+currentProfile.Xp)/100))) + 1

	isYesterday := isYesterday(currentProfile.LastActive)
	isToday := isToday(req.Date)
	var streak int
	if isYesterday && isToday {
		streak = currentProfile.Streak + 1
	} else {
		streak = 1
	}

	updatedProfile, err := s.levelRepo.UpdateLevel(ctx, &models.UserLevel{
		UserID:     usr.UserID,
		Xp:         req.Xp + currentProfile.Xp,
		Level:      int(level),
		Streak:     streak,
		LastActive: req.Date,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user level: %w", err)
	}

	return updatedProfile, nil
}

func (s *ProfileService) GetLevelService(ctx context.Context, usr models.UserIdentity) (*models.UserLevel, error) {
	profile, err := s.levelRepo.GetUserLevel(ctx, usr.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user level: %w", err)
	}

	return profile, nil
}

func (s *ProfileService) DeleteLevelService(ctx context.Context, usr models.UserIdentity) error {
	if err := s.levelRepo.DeleteUserLevel(ctx, usr.UserID); err != nil {
		return fmt.Errorf("failed to delete user level: %w", err)
	}
	return nil
}
