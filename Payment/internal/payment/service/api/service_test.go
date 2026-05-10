package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
	"github.com/sewaustav/Payment/internal/payment/service/api"
	"github.com/sewaustav/Payment/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func adminUser() models.UserIdentity {
	role := models.Admin
	return models.UserIdentity{
		UserID: 1,
		Role:   &role,
	}
}

func regularUser() models.UserIdentity {
	role := models.User
	return models.UserIdentity{
		UserID: 2,
		Role:   &role,
	}
}

// ─── RegisterUserService ──────────────────────────────────────────────────────

func TestRegisterUserService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	repo.On("InitSubscription", ctx, mock.MatchedBy(func(s *models.SubscriptionInfo) bool {
		return s.UserID == usr.UserID &&
			s.Subscription == models.NoSubscription &&
			s.CountOfRenewal == 0 &&
			!s.IsAutoRenew
	})).Return(&models.SubscriptionInfo{ID: 1}, nil)

	err := svc.RegisterUserService(ctx, usr)
	assert.NoError(t, err)
}

func TestRegisterUserService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()
	repoErr := errors.New("db error")

	repo.On("InitSubscription", ctx, mock.Anything).Return((*models.SubscriptionInfo)(nil), repoErr)

	err := svc.RegisterUserService(ctx, usr)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GetStatusService ─────────────────────────────────────────────────────────

func TestGetStatusService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	expected := dto.SubscriptionStatusDto{Status: 1}
	repo.On("GetSubscriptionStatus", ctx, usr.UserID).Return(expected, nil)

	result, err := svc.GetStatusService(ctx, usr)
	assert.NoError(t, err)
	assert.Equal(t, &expected, result)
}

func TestGetStatusService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()
	repoErr := errors.New("db error")

	repo.On("GetSubscriptionStatus", ctx, usr.UserID).Return(dto.SubscriptionStatusDto{}, repoErr)

	result, err := svc.GetStatusService(ctx, usr)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GetSubscriptionInfoService ──────────────────────────────────────────────

func TestGetSubscriptionInfoService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	expected := &models.SubscriptionInfo{ID: 1, UserID: usr.UserID}
	repo.On("GetUserSubscriptionInfo", ctx, usr.UserID).Return(expected, nil)

	result, err := svc.GetSubscriptionInfoService(ctx, usr)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetSubscriptionInfoService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()
	repoErr := errors.New("db error")

	repo.On("GetUserSubscriptionInfo", ctx, usr.UserID).Return((*models.SubscriptionInfo)(nil), repoErr)

	result, err := svc.GetSubscriptionInfoService(ctx, usr)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GetMyPaymentsService ─────────────────────────────────────────────────────

func TestGetMyPaymentsService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	expected := []models.PaymentInfo{{ID: 10}, {ID: 11}}
	// page=1, limit=10 → offset=0
	repo.On("GetUserPayments", ctx, usr.UserID, 10, 0).Return(expected, nil)

	result, err := svc.GetMyPaymentsService(ctx, usr, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetMyPaymentsService_LimitCapped(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	// limit=100 should be capped to 50; page=2 → offset = (2-1)*50 = 50
	repo.On("GetUserPayments", ctx, usr.UserID, 50, 50).Return([]models.PaymentInfo{}, nil)

	result, err := svc.GetMyPaymentsService(ctx, usr, 100, 2)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetMyPaymentsService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()
	repoErr := errors.New("db error")

	repo.On("GetUserPayments", ctx, usr.UserID, 10, 0).Return(([]models.PaymentInfo)(nil), repoErr)

	result, err := svc.GetMyPaymentsService(ctx, usr, 10, 1)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── UpdateSubscriptionInfoService ───────────────────────────────────────────

func TestUpdateSubscriptionInfoService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	subType := models.Basic
	isAutoRenew := true
	input := dto.UpdateSubscriptionInfoDto{
		Subscription: &subType,
		IsAutoRenew:  &isAutoRenew,
	}

	repo.On("UpdateSubscription", ctx, usr.UserID, mock.MatchedBy(func(d *dto.UpdateSubscriptionInfoDto) bool {
		return d.Subscription == input.Subscription &&
			d.IsAutoRenew == input.IsAutoRenew &&
			!d.IsRenew
	})).Return(nil)

	err := svc.UpdateSubscriptionInfoService(ctx, usr, input)
	assert.NoError(t, err)
}

func TestUpdateSubscriptionInfoService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()
	repoErr := errors.New("db error")

	repo.On("UpdateSubscription", ctx, usr.UserID, mock.Anything).Return(repoErr)

	err := svc.UpdateSubscriptionInfoService(ctx, usr, dto.UpdateSubscriptionInfoDto{})
	assert.ErrorIs(t, err, repoErr)
}

// ─── DeleteUserService ────────────────────────────────────────────────────────

func TestDeleteUserService_Success_Admin(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(99)

	repo.On("DeleteUser", ctx, targetID).Return(nil)

	err := svc.DeleteUserService(ctx, usr, targetID)
	assert.NoError(t, err)
}

func TestDeleteUserService_NonAdmin(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	err := svc.DeleteUserService(ctx, usr, 99)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not admin")
}

func TestDeleteUserService_NilRole_Allowed(t *testing.T) {
	// When Role == nil, the condition `usr.Role != nil && *usr.Role != models.Admin` is false,
	// so DeleteUser is called.
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := models.UserIdentity{UserID: 5, Role: nil}

	repo.On("DeleteUser", ctx, int64(99)).Return(nil)

	err := svc.DeleteUserService(ctx, usr, 99)
	assert.NoError(t, err)
}

func TestDeleteUserService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	repoErr := errors.New("db error")

	repo.On("DeleteUser", ctx, int64(99)).Return(repoErr)

	err := svc.DeleteUserService(ctx, usr, 99)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GiftSubscription ─────────────────────────────────────────────────────────

func TestGiftSubscription_NilRole_Error(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := models.UserIdentity{UserID: 1, Role: nil}

	err := svc.GiftSubscription(ctx, usr, 42)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "role is required")
}

func TestGiftSubscription_NonAdmin_Error(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := regularUser()

	err := svc.GiftSubscription(ctx, usr, 42)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not admin")
}

func TestGiftSubscription_UserAlreadyHasSubscription(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)

	sub := &models.SubscriptionInfo{ID: 1, UserID: targetID, Subscription: models.Basic}
	repo.On("GetUserSubscriptionInfo", ctx, targetID).Return(sub, nil)

	err := svc.GiftSubscription(ctx, usr, targetID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already has an active subscription")
}

func TestGiftSubscription_GetSubInfoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)
	repoErr := errors.New("db error")

	repo.On("GetUserSubscriptionInfo", ctx, targetID).Return((*models.SubscriptionInfo)(nil), repoErr)

	err := svc.GiftSubscription(ctx, usr, targetID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "internal")
}

func TestGiftSubscription_NewUser_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)

	// No existing subscription record
	repo.On("GetUserSubscriptionInfo", ctx, targetID).Return((*models.SubscriptionInfo)(nil), nil)

	mockTx := mocks.NewTx(t)
	repo.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)

	txRepo := mocks.NewPaymentRepo(t)
	repo.On("WithTx", mockTx).Return(txRepo)

	newSub := &models.SubscriptionInfo{ID: 10, UserID: targetID}
	txRepo.On("InitSubscription", ctx, mock.MatchedBy(func(s *models.SubscriptionInfo) bool {
		return s.UserID == targetID && s.Subscription == models.Basic
	})).Return(newSub, nil)

	txRepo.On("CreatePayment", ctx, mock.MatchedBy(func(p *models.PaymentInfo) bool {
		return p.UserID == targetID && p.Price == 0 && p.Status == "gift"
	})).Return(&models.PaymentInfo{ID: 1}, nil)

	mockTx.On("Commit").Return(nil)

	err := svc.GiftSubscription(ctx, usr, targetID)
	assert.NoError(t, err)
}

func TestGiftSubscription_ExistingUser_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)

	// Existing user with NO active subscription (Subscription == 0 == NoSubscription)
	existingSub := &models.SubscriptionInfo{ID: 5, UserID: targetID, Subscription: models.NoSubscription}
	repo.On("GetUserSubscriptionInfo", ctx, targetID).Return(existingSub, nil)

	mockTx := mocks.NewTx(t)
	repo.On("Begin", ctx).Return(mockTx, nil)
	mockTx.On("Rollback").Return(nil)

	txRepo := mocks.NewPaymentRepo(t)
	repo.On("WithTx", mockTx).Return(txRepo)

	txRepo.On("UpdateSubscription", ctx, targetID, mock.Anything).Return(nil)

	txRepo.On("CreatePayment", ctx, mock.MatchedBy(func(p *models.PaymentInfo) bool {
		return p.UserID == targetID && p.Price == 0 && p.Status == "gift"
	})).Return(&models.PaymentInfo{ID: 2}, nil)

	mockTx.On("Commit").Return(nil)

	err := svc.GiftSubscription(ctx, usr, targetID)
	assert.NoError(t, err)
}

// ─── GetUserProfileService ────────────────────────────────────────────────────

func TestGetUserProfileService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)

	expected := &models.SubscriptionInfo{ID: 1, UserID: targetID}
	repo.On("GetUserSubscriptionInfo", ctx, targetID).Return(expected, nil)

	result, err := svc.GetUserProfileService(ctx, usr, targetID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetUserProfileService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)
	repoErr := errors.New("db error")

	repo.On("GetUserSubscriptionInfo", ctx, targetID).Return((*models.SubscriptionInfo)(nil), repoErr)

	result, err := svc.GetUserProfileService(ctx, usr, targetID)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GetUsersPaymentsService ──────────────────────────────────────────────────

func TestGetUsersPaymentsService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)

	expected := []models.PaymentInfo{{ID: 1}, {ID: 2}}
	// limit=10, page=1 → offset=0
	repo.On("GetUserPayments", ctx, targetID, 10, 0).Return(expected, nil)

	result, err := svc.GetUsersPaymentsService(ctx, usr, targetID, 10, 1)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetUsersPaymentsService_LimitCapped(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)

	// limit=200 → capped to 50; page=3 → offset=100
	repo.On("GetUserPayments", ctx, targetID, 50, 100).Return([]models.PaymentInfo{}, nil)

	result, err := svc.GetUsersPaymentsService(ctx, usr, targetID, 200, 3)
	assert.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetUsersPaymentsService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	targetID := int64(42)
	repoErr := errors.New("db error")

	repo.On("GetUserPayments", ctx, targetID, 10, 0).Return(([]models.PaymentInfo)(nil), repoErr)

	result, err := svc.GetUsersPaymentsService(ctx, usr, targetID, 10, 1)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GetPaymentByTransactionIDService ────────────────────────────────────────

func TestGetPaymentByTransactionIDService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	txID := "txn-123"

	expected := &models.PaymentInfo{ID: 7, TransactionID: &txID}
	repo.On("GetPaymentByTransactionID", ctx, txID).Return(expected, nil)

	result, err := svc.GetPaymentByTransactionIDService(ctx, usr, txID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetPaymentByTransactionIDService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	txID := "txn-999"
	repoErr := errors.New("not found")

	repo.On("GetPaymentByTransactionID", ctx, txID).Return((*models.PaymentInfo)(nil), repoErr)

	result, err := svc.GetPaymentByTransactionIDService(ctx, usr, txID)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── GetPaymentByIDService ────────────────────────────────────────────────────

func TestGetPaymentByIDService_Success(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	payID := int64(55)

	expected := &models.PaymentInfo{ID: payID}
	repo.On("GetPaymentByID", ctx, payID).Return(expected, nil)

	result, err := svc.GetPaymentByIDService(ctx, usr, payID)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestGetPaymentByIDService_RepoError(t *testing.T) {
	repo := mocks.NewPaymentRepo(t)
	svc := service.NewPaymentService(repo)
	ctx := context.Background()
	usr := adminUser()
	payID := int64(55)
	repoErr := errors.New("not found")

	repo.On("GetPaymentByID", ctx, payID).Return((*models.PaymentInfo)(nil), repoErr)

	result, err := svc.GetPaymentByIDService(ctx, usr, payID)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
}

// ─── compile-time interface check ────────────────────────────────────────────

var _ service.PaymentApiService = (*service.PaymentApiCore)(nil)

// ensure time import is used (used implicitly via time.Now() inside the service itself;
// this blank assignment prevents an "imported and not used" error if time is referenced nowhere else)
var _ = time.Now
