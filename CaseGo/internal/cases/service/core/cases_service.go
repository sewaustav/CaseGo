package service

import (
	"context"
	"errors"

	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (s *CaseGoCoreService) GetCasesService(ctx context.Context, limit, page int, settings *dto.UserSettingsDto) ([]models.Case, error) {
	offset := (page - 1) * limit
	if settings.Topic != nil {
		cases, err := s.caseGoRepo.GetCasesByTopic(ctx, *settings.Topic, limit, offset)
		if err != nil {
			return nil, apperrors.NewInternal("failed to get cases by topic", err)
		}
		return cases, nil
	}

	if settings.Category != nil {
		cases, err := s.caseGoRepo.GetCasesByCategory(ctx, *settings.Category, limit, offset)
		if err != nil {
			return nil, apperrors.NewInternal("failed to get cases by category", err)
		}
		return cases, nil
	}

	cases, err := s.caseGoRepo.GetCases(ctx, limit, offset)
	if err != nil {
		return nil, apperrors.NewInternal("failed to get cases", err)
	}
	return cases, nil
}

func (s *CaseGoCoreService) GetCaseByIDService(ctx context.Context, caseID int64) (*models.Case, error) {
	caseModel, err := s.caseGoRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("case not found", err)
		}
		return nil, apperrors.NewInternal("failed to get case", err)
	}
	return caseModel, nil
}

func (s *CaseGoCoreService) CreateCaseService(ctx context.Context, caseDto *dto.NewCaseDto, user models.UserIdentity) (*models.Case, error) {
	if user.Role != models.Admin && user.Role != models.Creator {
		return nil, apperrors.NewForbidden("only admin and creator can create case", nil)
	}

	if caseDto.Prompt != nil {
		newCase, err := s.llmService.GenerateCase(ctx, *caseDto.Prompt)
		if err != nil {
			return nil, apperrors.NewInternal("failed to generate case via LLM", err)
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
		if errors.Is(err, apperrors.ErrAlreadyExists) {
			return nil, apperrors.NewConflict("case already exists", err)
		}
		return nil, apperrors.NewInternal("failed to create case", err)
	}
	return newCase, nil
}

func (s *CaseGoCoreService) DeleteCaseService(ctx context.Context, caseID int64, user models.UserIdentity) error {
	if user.Role != models.Admin {
		return apperrors.NewForbidden("only admin can delete case", nil)
	}
	err := s.caseGoRepo.DeleteCase(ctx, caseID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.NewNotFound("case not found", err)
		}
		return apperrors.NewInternal("failed to delete case", err)
	}
	return nil
}

func (s *CaseGoCoreService) PatchCaseService(ctx context.Context, caseID int64, caseDto *dto.NewCaseDto, user models.UserIdentity) (*models.Case, error) {
	if user.Role != models.Admin && user.Role != models.Creator {
		return nil, apperrors.NewForbidden("only admin and creator can patch case", nil)
	}

	if caseDto == nil {
		return nil, apperrors.NewBadRequest("case data is required", nil)
	}

	result, err := s.caseGoRepo.PatchCase(ctx, &models.Case{
		ID:            caseID,
		Topic:         *caseDto.Topic,
		Description:   *caseDto.Description,
		Category:      *caseDto.Category,
		IsGenerated:   false,
		FirstQuestion: *caseDto.FirstQuestion,
		Creator:       user.UserID,
	})
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.NewNotFound("case not found", err)
		}
		return nil, apperrors.NewInternal("failed to patch case", err)
	}
	return result, nil
}
