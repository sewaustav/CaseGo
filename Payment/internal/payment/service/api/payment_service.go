package service

import (
	"context"
	"fmt"

	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
	"github.com/sewaustav/Payment/internal/payment/repository"
)

type PaymentApiCore struct {
	repo repository.PaymentRepo
}

func (s *PaymentApiCore) GetStatusService(ctx context.Context, usr models.UserIdentity) (*dto.SubscriptionStatusDto, error) {
	status, err := s.repo.GetSubscriptionStatus(ctx, usr.UserID)
	if err != nil {
		return nil, err 
	}

	return &status, err
}

func (s *PaymentApiCore) GetSubscriptionInfoService(ctx context.Context, usr models.UserIdentity) (*models.SubscriptionInfo, error) {
	userSubscription, err := s.repo.GetUserSubscriptionInfo(ctx, usr.UserID)
	if err != nil {
		return nil, err
	}

	return userSubscription, nil
}

func (s *PaymentApiCore) GetMyPaymentsService(ctx context.Context, usr models.UserIdentity, limit, page int) ([]models.PaymentInfo, error) {
	offset := (page - 1) * limit 
	payments, err := s.repo.GetUserPayments(ctx, usr.UserID, limit, offset)
	if err != nil {
		return nil, err 
	}

	return payments, nil
	
}

func (s *PaymentApiCore) UpdateSubscriptionInfoService(ctx context.Context, usr models.UserIdentity, sub dto.UpadateSubcriptionInfoDto) error {
	changes := &dto.UpadateSubcriptionInfoDto{
		Subscription: sub.Subscription,
		IsAutoRenew: sub.IsAutoRenew,
		IsRenew: false,
	}
	if err := s.repo.UpdateSubscription(ctx, usr.UserID, changes); err != nil {
		return err
	}

	return nil
}

func (s *PaymentApiCore) DeleteUserService(ctx context.Context, usr models.UserIdentity) error {
	if usr.Role != nil && *usr.Role != models.Admin {
		return fmt.Errorf("user is not admin")
	}

	if err := s.repo.DeleteUser(ctx, usr.UserID); err != nil {
		return err 
	}

	return nil
}

// admins only 

func (s *PaymentApiCore) GetUserProfileService(ctx context.Context, usr models.UserIdentity, userID int64) (*models.SubscriptionInfo, error) {
	if usr.Role != nil && *usr.Role != models.Admin {
		return nil, fmt.Errorf("user is not admin")
	} 

	profile, err := s.repo.GetUserSubscriptionInfo(ctx, userID)
	if err != nil {
		return nil, err 
	}

	return profile, nil 
}

func (s *PaymentApiCore) GetUsersPaymentsService(ctx context.Context, usr models.UserIdentity, userID int64, limit, page int) ([]models.PaymentInfo, error) {
	if usr.Role != nil && *usr.Role != models.Admin {
		return nil, fmt.Errorf("user is not admin")
	} 

	offset := (page - 1) * limit

	history, err := s.repo.GetUserPayments(ctx, userID, limit, offset)
	if err != nil {
		return nil, err 
	}

	return history, err
}

func (s *PaymentApiCore) GetPaymentByTransactionIDService(ctx context.Context, usr models.UserIdentity, id string) (*models.PaymentInfo, error) {
	if usr.Role != nil && *usr.Role != models.Admin {
		return nil, fmt.Errorf("user is not admin")
	} 

	payment, err := s.repo.GetPaymentByTransactionID(ctx, id) 
	if err != nil {
		return nil, err 
	}

	return payment, nil 
}

func (s *PaymentApiCore) GetPaymentByIDService(ctx context.Context, usr models.UserIdentity, id int64) (*models.PaymentInfo, error) {
	if usr.Role != nil && *usr.Role != models.Admin {
		return nil, fmt.Errorf("user is not admin")
	} 

	payment, err := s.repo.GetPaymentByID(ctx, id) 
	if err != nil {
		return nil, err 
	}

	return payment, nil
}
