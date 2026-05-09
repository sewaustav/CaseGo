package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sewaustav/Payment/internal/payment/dto"
	"github.com/sewaustav/Payment/internal/payment/models"
	"github.com/sewaustav/Payment/internal/payment/repository"
)

type PaymentApiCore struct {
	repo repository.PaymentRepo
}

func NewPaymentService(repo repository.PaymentRepo) *PaymentApiCore {
	return &PaymentApiCore{
		repo: repo,
	}
}

func (s *PaymentApiCore) RegisterUser(ctx context.Context, usr models.UserIdentity) error {
	now := time.Now()
	_, err := s.repo.InitSubscription(ctx, &models.SubscriptionInfo{
		UserID: usr.UserID,
		Subscription: models.NoSubscription,
		CountOfRenewal: 0,
		IsAutoRenew: false,
		ExpiredAt: now,
	})
	if err != nil {
		return err
	}

	return nil
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
	if limit > 50 {
		limit = 50
	}
	offset := (page - 1) * limit 
	payments, err := s.repo.GetUserPayments(ctx, usr.UserID, limit, offset)
	if err != nil {
		return nil, err 
	}

	return payments, nil
	
}

func (s *PaymentApiCore) UpdateSubscriptionInfoService(ctx context.Context, usr models.UserIdentity, sub dto.UpdateSubscriptionInfoDto) error {
	changes := &dto.UpdateSubscriptionInfoDto{
		Subscription: sub.Subscription,
		IsAutoRenew: sub.IsAutoRenew,
		IsRenew: false,
	}
	if err := s.repo.UpdateSubscription(ctx, usr.UserID, changes); err != nil {
		return err
	}

	return nil
}

func (s *PaymentApiCore) DeleteUserService(ctx context.Context, usr models.UserIdentity, userID int64) error {
	if usr.Role != nil && *usr.Role != models.Admin {
		return fmt.Errorf("user is not admin")
	}

	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return err 
	}

	return nil
}

// admins only 
func (s *PaymentApiCore) GiftSubscription(ctx context.Context, usr models.UserIdentity, userID int64) error {
	if usr.Role == nil {
		return fmt.Errorf("role is required")
	}

	if *usr.Role != models.Admin {
		return fmt.Errorf("user is not admin")
	}

	subInfo, err := s.repo.GetUserSubscriptionInfo(ctx, userID) 
	if err != nil {
		return fmt.Errorf("internal %s", err)
	}

	if subInfo != nil && subInfo.Subscription > 0 {
		return fmt.Errorf("user already has an active subscription")
	}

	now := time.Now()

	tx, err := s.repo.Begin(ctx)
	if err != nil {
		return err 
	}

	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	var subID int64
	
	if subInfo == nil {
		
		sub, err := txRepo.InitSubscription(ctx, &models.SubscriptionInfo{
			UserID: userID,
			Subscription: models.Basic,
			CountOfRenewal: 1,
			IsAutoRenew: false,
			FirstPaymentDate: &now,
			LastPaymentDate: &now,
			ExpiredAt: now.AddDate(0, 1, 0),
		})

		if err != nil {
			return err 
		}

		subID = sub.ID
		
	} else {
		subType := models.Basic
		newSubDate := now.AddDate(0, 1, 0)
		err := txRepo.UpdateSubscription(ctx, userID, &dto.UpdateSubscriptionInfoDto{
			IsRenew:      false,
			Subscription: &subType, 
			ExpiredAt: &newSubDate,
		})
		if err != nil {
			return fmt.Errorf("failed to update sub: %w", err)
		}
		subID = subInfo.ID
	}

	_, err = txRepo.CreatePayment(ctx, &models.PaymentInfo{
		UserID:         userID,
		SubscriptionID: &subID,
		Price:          0,
		Currency:       "RUB",
		Status:         "gift",
		PaymentDate:    now,
	})
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil

}

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
	
	if limit > 50 {
		limit = 50
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
