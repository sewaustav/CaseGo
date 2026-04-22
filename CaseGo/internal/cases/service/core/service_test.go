package service

import (
	"context"
	"net/http"

	"testing"

	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	"github.com/sewaustav/CaseGoCore/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type txMock struct {
	committed bool
	rolled    bool
}

func (t *txMock) Commit() error {
	t.committed = true
	return nil
}

func (t *txMock) Rollback() error {
	t.rolled = true
	return nil
}

type noopGRPC struct{}

func (n noopGRPC) SendResults(ctx context.Context, msg models.Result) error {
	return nil
}

func TestGetCasesService_WithTopic(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	topic := "go"
	expected := []models.Case{{ID: 1}}

	caseRepo.On("GetCasesByTopic", ctx, topic, 10, 0).Return(expected, nil)

	got, err := svc.GetCasesService(ctx, 10, 1, &dto.UserSettingsDto{Topic: &topic})

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetCasesService_WithCategory(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	category := int32(2)
	expected := []models.Case{{ID: 10}}

	caseRepo.On("GetCasesByCategory", ctx, category, 5, 5).Return(expected, nil)

	got, err := svc.GetCasesService(ctx, 5, 2, &dto.UserSettingsDto{Category: &category})

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetCaseByIDService(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	expected := &models.Case{ID: 42}
	caseRepo.On("GetCaseByID", ctx, int64(42)).Return(expected, nil)

	got, err := svc.GetCaseByIDService(ctx, 42)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestStartDialogService(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	user := models.UserIdentity{UserID: 7}
	expectedCase := &models.Case{ID: 11}

	caseRepo.On("GetCaseByID", ctx, int64(11)).Return(expectedCase, nil)
	dialogRepo.On("StartDialog", ctx, &models.Dialog{UserID: 7, CaseID: 11}).Return(&models.Dialog{ID: 100}, nil)

	got, err := svc.StartDialogService(ctx, 11, user)

	require.NoError(t, err)
	assert.Equal(t, int64(100), got.DialogID)
	assert.Equal(t, int64(11), got.CaseID)
}

func TestHandleInteractionService(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	user := models.UserIdentity{UserID: 7}
	interaction := &dto.InteractionDto{
		DialogID: 1,
		Step:     3,
		Question: "q",
		Answer:   "a",
	}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(&models.Dialog{ID: 1, UserID: 7}, nil)
	caseRepo.On("GetCaseByID", ctx, mock.Anything).Return(&models.Case{ID: 0}, nil)
	redisClient.On("GetFullHistory", ctx, int64(1)).Return([]models.Interaction{}, nil)
	llm.On("GenerateResponse", 
    ctx, 
    mock.Anything, // для *models.Case
    mock.Anything, // для *models.Dialog
    mock.MatchedBy(func(history []models.Interaction) bool {
        return len(history) == 1 &&
            history[0].DialogID == 1 &&
            history[0].Step == 3 &&
            history[0].Question == "q" &&
            history[0].Answer == "a"
    }),
).Return(&dto.CaseDto{Model: "gpt"}, nil)
	redisClient.On("Push", ctx, mock.MatchedBy(func(inter *models.Interaction) bool {
		return inter.DialogID == 1 &&
			inter.Step == 3 &&
			inter.Question == "q" &&
			inter.Answer == "a"
	})).Return(nil)

	got, err := svc.HandleInteractionService(ctx, interaction, user)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.DialogID)
	assert.Equal(t, "q", got.Question)
	assert.Equal(t, "gpt", got.Model)
	if assert.NotNil(t, got.Step) {
		assert.Equal(t, int32(4), *got.Step)
	}
}

func TestHandleInteractionService_ForbiddenUser(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	user := models.UserIdentity{UserID: 7}
	interaction := &dto.InteractionDto{DialogID: 1}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(&models.Dialog{ID: 1, UserID: 8}, nil)

	got, err := svc.HandleInteractionService(ctx, interaction, user)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
	assert.Contains(t, appErr.Message, "not authorized")
}

func TestCompleteDialogService(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, noopGRPC{})

	user := models.UserIdentity{UserID: 7}
	dialog := &models.Dialog{ID: 1, UserID: 7, CaseID: 99}
	history := []models.Interaction{{DialogID: 1, Step: 1}}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(dialog, nil)
	redisClient.On("GetFullHistory", ctx, int64(1)).Return(history, nil)

	tx := mocks.NewTx(t)
	interactionRepo.On("Begin", ctx).Return(tx, nil)
	interactionRepo.On("WithTx", tx).Return(interactionRepo)

	interactionRepo.On("PutInteraction", ctx, &history[0]).Return(nil)
	tx.On("Commit").Return(nil)
	tx.On("Rollback").Return(nil)

	llm.On("AnalyzeCase", ctx, mock.MatchedBy(func(conv []models.Interaction) bool {
		return len(conv) == 1 &&
			conv[0].DialogID == 1 &&
			conv[0].Step == 1
	})).Return(&dto.Result{
		StepsCount: 1,
		SkillsRating: dto.Skills{
			Assertiveness:        1,
			Empathy:              2,
			ClarityCommunication: 3,
			Resistance:           4,
			Eloquence:            5,
			Initiative:           6,
		},
	}, nil)

	redisClient.On("Clear", ctx, int64(1)).Return(nil)

	got, err := svc.CompleteDialogService(ctx, 1, user)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int32(1), got.StepsCount)
}

func TestGetUsersDialogsService_OnlyOwnerOrAdmin(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	user := models.UserIdentity{UserID: 10, Role: models.User}
	got, err := svc.GetUsersDialogsService(ctx, user, 11, 10, 0)

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestGetUserDialogByIDService(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	grpcHandler := mocks.NewGRPCService(t)
	llm := mocks.NewLLM(t)

	svc := NewCaseGoCoreService(redisClient, caseRepo, dialogRepo, interactionRepo, llm, grpcHandler)

	user := models.UserIdentity{UserID: 7}
	dialog := &models.Dialog{ID: 1, UserID: 7}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(dialog, nil)
	redisClient.On("GetFullHistory", ctx, int64(1)).Return([]models.Interaction{{DialogID: 1}}, nil)

	got, err := svc.GetUserDialogByIDService(ctx, user, 1)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.Dialog.ID)
}
