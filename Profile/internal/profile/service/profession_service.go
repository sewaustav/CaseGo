package profileService

import (
	"context"
	"fmt"

	dto "github.com/YoungFlores/Case_Go/Profile/internal/profile/dto"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	apperr "github.com/YoungFlores/Case_Go/Profile/pkg/errors"
)

func (s *ProfileService) AddProfessionService(
	ctx context.Context,
	usr models.UserIdentity,
	req []dto.ProfessionDTO,
) ([]models.UserProfession, error) {

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin tx: %w", apperr.ErrInternal)
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	var professions []models.UserProfession

	for _, prof := range req {
		p, err := txRepo.AddProfession(ctx, &models.UserProfession{
			UserID:       usr.UserID,
			ProfessionID: prof.ProfessionID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add profession: %w", apperr.ErrInternal)
		}
		professions = append(professions, *p)

		parent, err := s.catRepo.GetParentOfCategory(ctx, prof.ProfessionID)
		if err != nil {
			return nil, fmt.Errorf("failed to add profession: %w", apperr.ErrInternal)
		}
		if parent != nil {
			p, err = txRepo.AddProfession(ctx, &models.UserProfession{
				UserID:       usr.UserID,
				ProfessionID: *parent,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to add profession: %w", apperr.ErrInternal)
			}
			if p != nil {
				professions = append(professions, *p)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", apperr.ErrInternal)
	}

	return professions, nil

}

func (s *ProfileService) EditProfessionCategoryService(ctx context.Context, usr models.UserIdentity, profession *models.UserProfession) (*models.UserProfession, error) {
	userID, err := s.repo.GetProfileIDByProfessionID(ctx, profession.ID)
	if err != nil {
		return nil, err
	}

	if userID != usr.UserID {
		return nil, apperr.ErrForbidden
	}

	return s.repo.EditProfession(ctx, profession)
}

func (s *ProfileService) DeleteProfessionService(ctx context.Context, id int64, usr models.UserIdentity) error {

	userID, err := s.repo.GetProfileIDByProfessionID(ctx, id)
	if err != nil {
		return err
	}

	if userID != usr.UserID {
		return apperr.ErrForbidden
	}

	return s.repo.DeleteProfession(ctx, id)

}

func (s *ProfileService) GetProfessionsService(ctx context.Context, usr models.UserIdentity) ([]models.UserProfession, error) {
	listProfessions, err := s.repo.GetAllProfessions(ctx, usr.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get professions: %w", err)
	}

	return listProfessions, nil
}
