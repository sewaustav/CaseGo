package subscription_repo

import (
	"context"

	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
)

type SubscriptionInfo interface {
	PushSubInfo(ctx context.Context, userID int64, sub dto.SubscriptionStatusDto) error
	GetSubInfo(ctx context.Context, userID int64) (*dto.SubscriptionStatusDto, error)
	Close() error
}