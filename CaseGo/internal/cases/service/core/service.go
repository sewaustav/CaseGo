package service

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cache"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/handlers/grpc"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	"github.com/sewaustav/CaseGoCore/internal/cases/repository"
	"github.com/sewaustav/CaseGoCore/internal/cases/service/llm_service"
)

type CaseGoService interface {
	StartDialogService(ctx context.Context, caseID int64, user models.UserIdentity) (*dto.StartDialogResponse, error)
	HandleInteractionService(ctx context.Context, interaction *dto.InteractionDto, user models.UserIdentity) (*dto.CaseDto, error)
	CompleteDialogService(ctx context.Context, dialogID int64, user models.UserIdentity) (*dto.Result, error)

	GetCasesService(ctx context.Context, limit, page int, settings *dto.UserSettingsDto) ([]models.Case, error)
	GetCaseByIDService(ctx context.Context, caseID int64) (*models.Case, error)

	// admin only
	CreateCaseService(ctx context.Context, caseDto *dto.NewCaseDto, identity models.UserIdentity) (*models.Case, error)
	DeleteCaseService(ctx context.Context, caseID int64, user models.UserIdentity) error
	PatchCaseService(ctx context.Context, caseID int64, caseDto *dto.NewCaseDto, user models.UserIdentity) (*models.Case, error)

	//users
	GetUsersDialogsService(ctx context.Context, user models.UserIdentity, userID int64, limit, offset int) ([]models.Conversation, error)
	GetUserDialogByIDService(ctx context.Context, user models.UserIdentity, dialogID int64) (*models.Conversation, error)
}

type CaseGoCoreService struct {
	redisClient     cache.Interactor
	caseGoRepo      repository.CaseRepo
	dialogRepo      repository.DialogRepo
	interactionRepo repository.Interaction
	llmService      llm_service.LLM
	grpcHandler     grpc.GRPCService
}

func NewCaseGoCoreService(redisClient cache.Interactor, caseGoRepo repository.CaseRepo, dialogRepo repository.DialogRepo, interactionRepo repository.Interaction, llm_service llm_service.LLM, grpc grpc.GRPCService) *CaseGoCoreService {
	return &CaseGoCoreService{
		redisClient:     redisClient,
		caseGoRepo:      caseGoRepo,
		dialogRepo:      dialogRepo,
		interactionRepo: interactionRepo,
		llmService:      llm_service,
		grpcHandler:     grpc,
	}
}
