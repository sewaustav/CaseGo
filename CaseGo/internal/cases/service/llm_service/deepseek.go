package llm_service

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

func (l *LLMService) GenerateCase(ctx context.Context, description string) (*models.Case, error) {
	return models.Case{}, nil
}

func (l *LLMService) GenerateResponse(ctx context.Context, history []models.Interaction) (*dto.CaseDto, error) {
	return dto.CaseDto{}, nil
}

func (l *LLMService) AnalyzeCase(ctx context.Context, conv []models.Interaction) (*dto.Result, error) {
	return dto.Result{}, nil
}
