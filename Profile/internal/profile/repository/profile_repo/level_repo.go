package profilerepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
)

type LevelRepo interface {
	CreateLevel(ctx context.Context, userLevel *models.UserLevel) (*models.UserLevel, error)
	UpdateLevel(ctx context.Context, userLevel *models.UserLevel) (*models.UserLevel, error)
	GetUserLevel(ctx context.Context, userID int64) (*models.UserLevel, error)
	DeleteUserLevel(ctx context.Context, userID int64) error
}

type PostgresLevelRepo struct {
	db *sql.DB
}

func NewPostgresLevelRepo(db *sql.DB) *PostgresLevelRepo {
	return &PostgresLevelRepo{
		db: db,
	}
}

func (r *PostgresLevelRepo) CreateLevel(ctx context.Context, userLevel *models.UserLevel) (*models.UserLevel, error) {
	now := time.Now()
	query, args, err := psql.Insert("levels").
		Columns("user_id", "xp", "level", "streak", "last_active").
		Values(userLevel.UserID, userLevel.Xp, userLevel.Level, userLevel.Streak, now).
		Suffix("RETURNING id, user_id, xp, level, streak, last_active").ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&userLevel.ID, &userLevel.UserID, &userLevel.Xp, &userLevel.Level, &userLevel.Streak, &userLevel.LastActive)
	if err != nil {
		return nil, err
	}
	return userLevel, nil
}

func (r *PostgresLevelRepo) UpdateLevel(ctx context.Context, userLevel *models.UserLevel) (*models.UserLevel, error) {
	now := time.Now()
	query, args, err := psql.Update("levels").
		Set("xp", userLevel.Xp).
		Set("level", userLevel.Level).
		Set("streak", userLevel.Streak).
		Set("last_active", now).
		Where(sq.Eq{"user_id": userLevel.UserID}).
		Suffix("RETURNING id, user_id, xp, level, streak, last_active").
		ToSql()
	if err != nil {
		return nil, err
	}
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&userLevel.ID, &userLevel.UserID, &userLevel.Xp, &userLevel.Level, &userLevel.Streak, &userLevel.LastActive)
	if err != nil {
		return nil, err
	}
	return userLevel, nil
}

func (r *PostgresLevelRepo) GetUserLevel(ctx context.Context, userID int64) (*models.UserLevel, error) {
	query, args, err := psql.Select("id", "user_id", "xp", "level", "streak", "last_active").From("levels").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, fmt.Errorf("build get user level query: %w", err)
	}
	row := r.db.QueryRowContext(ctx, query, args...)
	var level models.UserLevel
	err = row.Scan(&level.ID, &level.UserID, &level.Xp, &level.Level, &level.Streak, &level.LastActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &repoerr.RepoError{
				Field: "user_id",
				Err:   repoerr.ErrNotFound,
			}
		}
		return nil, fmt.Errorf("scan user level: %w", err)
	}
	return &level, nil
}

func (r *PostgresLevelRepo) DeleteUserLevel(ctx context.Context, userID int64) error {
	query, args, err := psql.Delete("levels").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return fmt.Errorf("build delete user level query: %w", err)
	}
	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete user level: %w", err)
	}
	return nil
}
