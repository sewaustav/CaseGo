package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

const ttl = 1*time.Hour + 10*time.Minute

type redisRepo struct {
	client *redis.Client
}

func New(addr, password string, db int) (Interactor, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	return &redisRepo{client: rdb}, nil
}

func (r *redisRepo) makeKey(id int64) string {
	return fmt.Sprintf("dialog:%d", id)
}

func (r *redisRepo) Push(ctx context.Context, inter *models.Interaction) error {
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

func (r *redisRepo) GetFullHistory(ctx context.Context, dialogID int64) ([]models.Interaction, error) {
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

func (r *redisRepo) DeleteLast(ctx context.Context, dialogID int64) error {
	return r.client.RPop(ctx, r.makeKey(dialogID)).Err()
}

func (r *redisRepo) Clear(ctx context.Context, dialogID int64) error {
	return r.client.Del(ctx, r.makeKey(dialogID)).Err()
}
