package service

import (
	"context"
	"errors"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (s *CaseGoCoreService) GetCasesService(ctx context.Context, limit, page int, settings *dto.UserSettingsDto) ([]models.Case, error) {
	offset := (page - 1) * limit
	if settings.Topic != nil {
		сases, err := s.caseGoRepo.GetCasesByTopic(ctx, *settings.Topic, limit, offset)
		if err != nil {
			return nil, err
		}
		return сases, nil
	} else if settings.Category != nil {
		cases, err := s.caseGoRepo.GetCasesByCategory(ctx, *settings.Category, limit, offset)
		if err != nil {
			return nil, err
		}
		return cases, nil
	}
	cases, err := s.caseGoRepo.GetCases(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return cases, nil
}

func (s *CaseGoCoreService) GetCaseByIDService(ctx context.Context, caseID int64) (*models.Case, error) {
	caseModel, err := s.caseGoRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return nil, err
	}
	return caseModel, nil
}

func (s *CaseGoCoreService) CreateCaseService(ctx context.Context, caseDto *dto.NewCaseDto, user models.UserIdentity) (*models.Case, error) {
	if user.Role == models.Admin || user.Role == models.Creator {
		if caseDto.Prompt != nil {
			newCase, err := s.llmService.GenerateCase(ctx, *caseDto.Prompt)
			if err != nil {
				return nil, err
			}
			return newCase, nil
		}
		newCase, err := s.caseGoRepo.CreateCase(ctx, &models.Case{
			Topic:         *caseDto.Topic,
			Description:   *caseDto.Description,
			Category:      *caseDto.Category,
			IsGenerated:   false,
			FirstQuestion: *caseDto.FirstQuestion,
			Creator:       user.UserID,
		})
		if err != nil {
			return nil, err
		}
		return newCase, nil

	}

	return nil, errors.New("only admin and creator can create case")

}

func (s *CaseGoCoreService) DeleteCaseService(ctx context.Context, caseID int64, user models.UserIdentity) error {
	if user.Role != models.Admin {
		return errors.New("only admin can delete case")
	}
	return s.caseGoRepo.DeleteCase(ctx, caseID)
}

func (s *CaseGoCoreService) PatchCaseService(ctx context.Context, caseID int64, caseDto *dto.NewCaseDto, user models.UserIdentity) (*models.Case, error) {
	if user.Role != models.Admin && user.Role != models.Creator {
		return nil, errors.New("only admin and creator can patch case")
	}
	if caseDto == nil {
		return nil, errors.New("caseDto is required")
	}
	return s.caseGoRepo.PatchCase(ctx, &models.Case{
		ID:            caseID,
		Topic:         *caseDto.Topic,
		Description:   *caseDto.Description,
		Category:      *caseDto.Category,
		IsGenerated:   false,
		FirstQuestion: *caseDto.FirstQuestion,
		Creator:       user.UserID,
	})
}
