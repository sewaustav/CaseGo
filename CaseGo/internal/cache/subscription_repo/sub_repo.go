package subscription_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sewaustav/CaseGoCore/internal/cases/dto"
)

const ttl = 5 * time.Minute

type repo struct {
	client *redis.Client
}

func New(client *redis.Client) SubscriptionInfo {
	return &repo{
		client: client,
	}
}

func (r *repo) makeKey(userID int64) string {
	return fmt.Sprintf("subscription:%d", userID)
}

func (r *repo) PushSubInfo(
	ctx context.Context, userID int64,
	sub dto.SubscriptionStatusDto,
) error {
	key := r.makeKey(userID)

	data, err := json.Marshal(sub)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, key, data, ttl).Err()
}

func (r *repo) GetSubInfo(
	ctx context.Context,
	userID int64,
) (*dto.SubscriptionStatusDto, error) {
	key := r.makeKey(userID)

	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, err
	}

	var sub dto.SubscriptionStatusDto

	if err := json.Unmarshal([]byte(val), &sub); err != nil {
		return nil, err
	}

	return &sub, nil
}