package service

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
	"github.com/sewaustav/CaseGoCore/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type noopGRPC struct{}

func (n noopGRPC) SendResults(ctx context.Context, msg models.Result) error {
	return nil
}

// activeSub returns a valid subscription that won't expire soon.
func activeSub() *dto.SubscriptionStatusDto {
	return &dto.SubscriptionStatusDto{
		Status:    1,
		ExpiredAt: time.Now().Add(2 * time.Hour),
	}
}

// newSvc is a helper that wires up all mocks and returns the service.
func newSvc(t *testing.T) (
	*CaseGoCoreService,
	*mocks.CaseRepo,
	*mocks.DialogRepo,
	*mocks.Interaction,
	*mocks.Interactor,
	*mocks.SubscriptionInfo,
	*mocks.GRPCService,
	*mocks.PaymentGrpcClient,
	*mocks.LLM,
) {
	t.Helper()
	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	subRedis := mocks.NewSubscriptionInfo(t)
	grpcHandler := mocks.NewGRPCService(t)
	paymentCheck := mocks.NewPaymentGrpcClient(t)
	llm := mocks.NewLLM(t)
	grpcLevel := mocks.NewLevelGrpcClient(t)

	svc := NewCaseGoCoreService(
		redisClient, subRedis,
		caseRepo, dialogRepo, interactionRepo,
		llm, grpcHandler, paymentCheck, grpcLevel,
	)
	return svc, caseRepo, dialogRepo, interactionRepo, redisClient, subRedis, grpcHandler, paymentCheck, llm
}

// stubSub стабит получение подписки через кэш (быстрый путь).
func stubSub(subRedis *mocks.SubscriptionInfo, userID int64) {
	subRedis.On("GetSubInfo", mock.Anything, userID).Return(activeSub(), nil)
}

// ────────────────────────────────────────────────────────────────────────────
// GetCasesService
// ────────────────────────────────────────────────────────────────────────────

func TestGetCasesService_WithTopic(t *testing.T) {
	ctx := context.Background()
	svc, caseRepo, _, _, _, _, _, _, _ := newSvc(t)

	topic := "go"
	expected := []models.Case{{ID: 1}}
	caseRepo.On("GetCasesByTopic", ctx, topic, 10, 0).Return(expected, nil)

	got, err := svc.GetCasesService(ctx, 10, 1, &dto.UserSettingsDto{Topic: &topic})

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestGetCasesService_WithCategory(t *testing.T) {
	ctx := context.Background()
	svc, caseRepo, _, _, _, _, _, _, _ := newSvc(t)

	category := int32(2)
	expected := []models.Case{{ID: 10}}
	caseRepo.On("GetCasesByCategory", ctx, category, 5, 5).Return(expected, nil)

	got, err := svc.GetCasesService(ctx, 5, 2, &dto.UserSettingsDto{Category: &category})

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

// ────────────────────────────────────────────────────────────────────────────
// GetCaseByIDService
// ────────────────────────────────────────────────────────────────────────────

func TestGetCaseByIDService(t *testing.T) {
	ctx := context.Background()
	svc, caseRepo, _, _, _, _, _, _, _ := newSvc(t)

	expected := &models.Case{ID: 42}
	caseRepo.On("GetCaseByID", ctx, int64(42)).Return(expected, nil)

	got, err := svc.GetCaseByIDService(ctx, 42)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

// ────────────────────────────────────────────────────────────────────────────
// StartDialogService
// ────────────────────────────────────────────────────────────────────────────

func TestStartDialogService(t *testing.T) {
	ctx := context.Background()
	svc, caseRepo, dialogRepo, _, _, subRedis, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	stubSub(subRedis, user.UserID)

	expectedCase := &models.Case{ID: 11}
	caseRepo.On("GetCaseByID", ctx, int64(11)).Return(expectedCase, nil)
	dialogRepo.On("StartDialog", ctx, &models.Dialog{UserID: 7, CaseID: 11}).
		Return(&models.Dialog{ID: 100}, nil)

	got, err := svc.StartDialogService(ctx, 11, user)

	require.NoError(t, err)
	assert.Equal(t, int64(100), got.DialogID)
	assert.Equal(t, int64(11), got.CaseID)
}

func TestStartDialogService_InactiveSubscription(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _, subRedis, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	subRedis.On("GetSubInfo", mock.Anything, user.UserID).Return(&dto.SubscriptionStatusDto{
		Status:    0,
		ExpiredAt: time.Now().Add(2 * time.Hour),
	}, nil)

	got, err := svc.StartDialogService(ctx, 11, user)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

func TestStartDialogService_ExpiredSubscription(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _, subRedis, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	subRedis.On("GetSubInfo", mock.Anything, user.UserID).Return(&dto.SubscriptionStatusDto{
		Status:    1,
		ExpiredAt: time.Now().Add(-1 * time.Hour),
	}, nil)

	got, err := svc.StartDialogService(ctx, 11, user)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

// ────────────────────────────────────────────────────────────────────────────
// HandleInteractionService
// ────────────────────────────────────────────────────────────────────────────

func TestHandleInteractionService(t *testing.T) {
	ctx := context.Background()
	svc, caseRepo, dialogRepo, _, redisClient, subRedis, _, _, llm := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	stubSub(subRedis, user.UserID)

	interaction := &dto.InteractionDto{
		DialogID: 1,
		Step:     3,
		Question: "q",
		Answer:   "a",
	}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).
		Return(&models.Dialog{ID: 1, UserID: 7}, nil)
	caseRepo.On("GetCaseByID", ctx, mock.Anything).
		Return(&models.Case{ID: 0}, nil)
	redisClient.On("GetFullHistory", ctx, int64(1)).
		Return([]models.Interaction{}, nil)

	nextStep := int32(4)
	llm.On("GenerateResponse",
		ctx,
		mock.Anything,
		mock.Anything,
		mock.MatchedBy(func(history []models.Interaction) bool {
			return len(history) == 1 &&
				history[0].DialogID == 1 &&
				history[0].Step == 3 &&
				history[0].Question == "q" &&
				history[0].Answer == "a"
		}),
	).Return(&dto.CaseDto{Model: "gpt", Question: "next question", Step: &nextStep}, nil)

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
	assert.Equal(t, "next question", got.Question)
	assert.Equal(t, "gpt", got.Model)
	if assert.NotNil(t, got.Step) {
		assert.Equal(t, int32(4), *got.Step)
	}
}

func TestHandleInteractionService_ForbiddenUser(t *testing.T) {
	ctx := context.Background()
	svc, _, dialogRepo, _, _, subRedis, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	stubSub(subRedis, user.UserID)

	interaction := &dto.InteractionDto{DialogID: 1}
	dialogRepo.On("GetDialogByID", ctx, int64(1)).
		Return(&models.Dialog{ID: 1, UserID: 8}, nil)

	got, err := svc.HandleInteractionService(ctx, interaction, user)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
	assert.Contains(t, appErr.Message, "not authorized")
}

func TestHandleInteractionService_NilInteraction(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _, subRedis, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	stubSub(subRedis, user.UserID)

	got, err := svc.HandleInteractionService(ctx, nil, user)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusBadRequest, appErr.Code)
}

// ────────────────────────────────────────────────────────────────────────────
// CompleteDialogService
// ────────────────────────────────────────────────────────────────────────────

func TestCompleteDialogService(t *testing.T) {
	ctx := context.Background()

	caseRepo := mocks.NewCaseRepo(t)
	dialogRepo := mocks.NewDialogRepo(t)
	interactionRepo := mocks.NewInteraction(t)
	redisClient := mocks.NewInteractor(t)
	subRedis := mocks.NewSubscriptionInfo(t)
	paymentCheck := mocks.NewPaymentGrpcClient(t)
	llm := mocks.NewLLM(t)
	grpcLevel := mocks.NewLevelGrpcClient(t)

	svc := NewCaseGoCoreService(
		redisClient, subRedis,
		caseRepo, dialogRepo, interactionRepo,
		llm, noopGRPC{}, paymentCheck, grpcLevel,
	)

	user := models.UserIdentity{UserID: 7}
	dialog := &models.Dialog{ID: 1, UserID: 7, CaseID: 99}
	history := []models.Interaction{{DialogID: 1, Step: 1}}
	xp := int32(10)

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(dialog, nil)
	caseRepo.On("GetCaseByID", ctx, int64(99)).Return(&models.Case{ID: 99, Xp: &xp}, nil)
	redisClient.On("GetFullHistory", ctx, int64(1)).Return(history, nil)

	tx := mocks.NewTx(t)
	interactionRepo.On("Begin", ctx).Return(tx, nil)
	interactionRepo.On("WithTx", tx).Return(interactionRepo)
	interactionRepo.On("PutInteraction", ctx, &history[0]).Return(nil)
	tx.On("Commit").Return(nil)
	tx.On("Rollback").Return(nil)

	llm.On("AnalyzeCase", ctx, mock.MatchedBy(func(conv []models.Interaction) bool {
		return len(conv) == 1 && conv[0].DialogID == 1 && conv[0].Step == 1
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

	grpcLevel.On("LevelDoneGrpcHandler", mock.Anything, user, mock.Anything).
		Return(&dto.UserLevelInfo{Level: 1, Xp: xp, IsLevelUp: false}, nil)

	got, err := svc.CompleteDialogService(ctx, 1, user)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int32(1), got.Result.StepsCount)
}

// ────────────────────────────────────────────────────────────────────────────
// GetUsersDialogsService
// ────────────────────────────────────────────────────────────────────────────

func TestGetUsersDialogsService_OnlyOwnerOrAdmin(t *testing.T) {
	ctx := context.Background()
	svc, _, _, _, _, _, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 10, Role: models.User}
	got, err := svc.GetUsersDialogsService(ctx, user, 11, 10, 0)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}

func TestGetUsersDialogsService_AdminCanViewOthers(t *testing.T) {
	ctx := context.Background()
	svc, _, dialogRepo, interactionRepo, redisClient, _, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 1, Role: models.Admin}
	dialogs := []models.Dialog{{ID: 5, UserID: 99}}
	history := []models.Interaction{{DialogID: 5}}

	dialogRepo.On("GetUserDialogs", ctx, int64(99), 10, 0).Return(dialogs, nil)
	redisClient.On("GetFullHistory", ctx, int64(5)).Return(history, nil)
	_ = interactionRepo

	got, err := svc.GetUsersDialogsService(ctx, user, 99, 10, 0)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, int64(5), got[0].Dialog.ID)
}

// ────────────────────────────────────────────────────────────────────────────
// GetUserDialogByIDService
// ────────────────────────────────────────────────────────────────────────────

func TestGetUserDialogByIDService(t *testing.T) {
	ctx := context.Background()
	svc, _, dialogRepo, _, redisClient, _, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7}
	dialog := &models.Dialog{ID: 1, UserID: 7}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(dialog, nil)
	redisClient.On("GetFullHistory", ctx, int64(1)).
		Return([]models.Interaction{{DialogID: 1}}, nil)

	got, err := svc.GetUserDialogByIDService(ctx, user, 1)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, int64(1), got.Dialog.ID)
}

func TestGetUserDialogByIDService_ForbiddenUser(t *testing.T) {
	ctx := context.Background()
	svc, _, dialogRepo, _, _, _, _, _, _ := newSvc(t)

	user := models.UserIdentity{UserID: 7, Role: models.User}
	dialog := &models.Dialog{ID: 1, UserID: 99}

	dialogRepo.On("GetDialogByID", ctx, int64(1)).Return(dialog, nil)

	got, err := svc.GetUserDialogByIDService(ctx, user, 1)

	require.Error(t, err)
	assert.Nil(t, got)

	var appErr *apperrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, http.StatusForbidden, appErr.Code)
}
