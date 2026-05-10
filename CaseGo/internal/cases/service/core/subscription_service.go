package service

import (
	"context"
	"time"

	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
)

func (s *CaseGoCoreService) getSubscriptionInfo(ctx context.Context, userID int64) (*dto.SubscriptionStatusDto, error) {
	subInfo, err := s.subRedisClient.GetSubInfo(ctx, userID)
	if err == nil && subInfo != nil {
		return subInfo, nil
	}

	subInfo, err = s.paymentCheck.CheckStatus(ctx, userID)
	if err != nil {
		return nil, err 
	}

	_ = s.subRedisClient.PushSubInfo(ctx, userID, *subInfo)

	return subInfo, nil
}

func (s *CaseGoCoreService) isSubscriptionValid(ctx context.Context, info *dto.SubscriptionStatusDto, rem float64) (bool, error) {
	if info.Status < 1 {
		return false, apperrors.NewForbidden("subscription is not active", nil) 
	}

	now := time.Now()

	if info.ExpiredAt.Before(now) {
		return false, apperrors.NewForbidden("subscription is out", nil)
	}

	remaining := info.ExpiredAt.Sub(now)
	if remaining.Minutes() < rem {
		return false, apperrors.NewForbidden("subsription is gone less then 10 minutes", nil)
	}

	return true, nil
}