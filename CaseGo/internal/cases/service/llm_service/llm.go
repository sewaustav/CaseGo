package llm_service

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type LLM interface {
	GenerateCase(ctx context.Context, topic, description string) (models.Case, error)
	GenerateResponse(ctx context.Context, history []models.Interaction) (dto.CaseDto, error)
	AnalyzeCase(ctx context.Context, conv []models.Interaction) (dto.Result, error)
}

type LLMService struct {
	URL string
}

func NewLLMService(url string) *LLMService {
	return &LLMService{
		URL: url,
	}
}
