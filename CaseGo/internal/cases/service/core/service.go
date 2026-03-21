package service

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cache"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	"github.com/sewaustav/CaseGoCore/internal/cases/repository"
	"github.com/sewaustav/CaseGoCore/internal/cases/service/llm_service"
)

type CaseGoService interface {
	StartDialog(ctx context.Context, caseID int64, user models.UserIdentity) (*models.Case, error)
	HandleInteraction(ctx context.Context, interaction *dto.InteractionDto) (*dto.CaseDto, error)
}

type CaseGoCoreService struct {
	redisClient     cache.Interactor
	caseGoRepo      repository.CaseRepo
	dialogRepo      repository.DialogRepo
	interactionRepo repository.Interaction
	llmService      llm_service.LLM
}

func NewCaseGoCoreService(redisClient cache.Interactor, caseGoRepo repository.CaseRepo, dialogRepo repository.DialogRepo, interactionRepo repository.Interaction, llm_service llm_service.LLM) *CaseGoCoreService {
	return &CaseGoCoreService{
		redisClient:     redisClient,
		caseGoRepo:      caseGoRepo,
		dialogRepo:      dialogRepo,
		interactionRepo: interactionRepo,
		llmService:      llm_service,
	}
}
