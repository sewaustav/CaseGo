package history_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

const ttl = 1*time.Hour + 10*time.Minute

type repo struct {
	client *redis.Client
}

func New(client *redis.Client) Interactor {
	return &repo{
		client: client,
	}
}

func (r *repo) makeKey(dialogID int64) string {
	return fmt.Sprintf("dialog:%d", dialogID)
}

func (r *repo) Push(ctx context.Context, inter *models.Interaction) error {
	key := r.makeKey(inter.DialogID)

	data, err := json.Marshal(inter)
	if err != nil {
		return err
	}

	pipe := r.client.Pipeline()

	pipe.RPush(ctx, key, data)
	pipe.Expire(ctx, key, ttl)

	_, err = pipe.Exec(ctx)
	return err
}

func (r *repo) GetFullHistory(
	ctx context.Context,
	dialogID int64,
) ([]models.Interaction, error) {
	key := r.makeKey(dialogID)

	vals, err := r.client.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	history := make([]models.Interaction, len(vals))

	for i, v := range vals {
		if err := json.Unmarshal([]byte(v), &history[i]); err != nil {
			return nil, err
		}
	}

	return history, nil
}

func (r *repo) DeleteLast(ctx context.Context, dialogID int64) error {
	return r.client.RPop(ctx, r.makeKey(dialogID)).Err()
}

func (r *repo) Clear(ctx context.Context, dialogID int64) error {
	return r.client.Del(ctx, r.makeKey(dialogID)).Err()
}
