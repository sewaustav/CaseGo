package service

import (
	"context"
	"time"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (s *CaseGoCoreService) StartDialog(ctx context.Context, caseID int64, user models.UserIdentity) (*models.Case, error) {
	caseModel, err := s.caseGoRepo.GetCaseByID(ctx, caseID)
	if err != nil {
		return nil, err
	}

	return caseModel, nil
}

func (s *CaseGoCoreService) HandleInteraction(ctx context.Context, interaction *dto.InteractionDto) (*dto.CaseDto, error) {
	interactionModel := &models.Interaction{
		DialogID:  interaction.DialogID,
		Step:      interaction.Step,
		Question:  interaction.Question,
		Answer:    interaction.Answer,
		CreatedAt: time.Now(),
	}

	history, err := s.redisClient.GetFullHistory(ctx, interaction.DialogID)
	if err != nil {
		return nil, err
	}

	history = append(history, *interactionModel)

	llmResponse, err := s.llmService.GenerateResponse(ctx, history)
	if err != nil {
		return nil, err
	}

	if err := s.redisClient.Push(ctx, interactionModel); err != nil {
		return nil, err
	}
	return &dto.CaseDto{
		DialogID: interaction.DialogID,
		Question: interaction.Question,
		Model:    llmResponse.Model,
		Step:     new(interaction.Step + 1),
	}, nil
}
