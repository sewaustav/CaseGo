package service

import (
	"context"

	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
)

type PaymentApiService interface {
	GetStatusService(ctx context.Context, usr models.UserIdentity) (*dto.SubscriptionStatusDto, error)
	GetSubscriptionInfoService(ctx context.Context, usr models.UserIdentity) (*models.SubscriptionInfo, error)
	GetMyPaymentsService(ctx context.Context, usr models.UserIdentity, limit, page int) ([]models.PaymentInfo, error)
	
	UpdateSubscriptionInfoService(ctx context.Context, usr models.UserIdentity, sub dto.UpadateSubcriptionInfoDto) error
	
	DeleteUserService(ctx context.Context, usr models.UserIdentity) error

	// for admin
	GetUserProfileService(ctx context.Context, usr models.UserIdentity, userID int64) (*models.SubscriptionInfo, error)
	GetUsersPaymentsService(ctx context.Context, usr models.UserIdentity, userID int64, limit, page int) ([]models.PaymentInfo, error)
	GetPaymentByTransactionIDService(ctx context.Context, usr models.UserIdentity, id string) (*models.PaymentInfo, error)
	GetPaymentByIDService(ctx context.Context, usr models.UserIdentity, id int64) (*models.PaymentInfo, error)
	
}