package service

import (
	"context"

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