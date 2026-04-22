package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/sewaustav/CaseGoCore/apperrors"
	"github.com/sewaustav/CaseGoCore/internal/cases/models"
)

type PostgresDialogRepo struct {
	db *sql.DB
}

func NewPostgresDialogRepo(db *sql.DB) *PostgresDialogRepo {
	return &PostgresDialogRepo{db: db}
}

func (r *PostgresDialogRepo) StartDialog(ctx context.Context, dialog *models.Dialog) (*models.Dialog, error) {
	query := psql.
		Insert("dialogs").
		Columns("case_id", "user_id", "model_name").
		Values(dialog.CaseID, dialog.UserID, dialog.ModelName).
		Suffix("RETURNING id, started_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("start dialog: build query: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&dialog.ID, &dialog.StartedAt); err != nil {
		return nil, fmt.Errorf("start dialog: exec: %w", err)
	}

	return dialog, nil
}

func (r *PostgresDialogRepo) GetDialogByID(ctx context.Context, dialogID int64) (*models.Dialog, error) {
	query := psql.
		Select("id", "case_id", "user_id", "model_name", "started_at", "ended_at").
		From("dialogs").
		Where(sq.Eq{"id": dialogID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get dialog by id: build query: %w", err)
	}

	var d models.Dialog
	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&d.ID,
		&d.CaseID,
		&d.UserID,
		&d.ModelName,
		&d.StartedAt,
		&d.EndedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("dialog id=%d: %w", dialogID, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get dialog by id: scan: %w", err)
	}

	return &d, nil
}

func (r *PostgresDialogRepo) GetUserDialogs(ctx context.Context, userID int64, limit, offset int) ([]models.Dialog, error) {
	query := psql.
		Select("id", "case_id", "user_id", "model_name", "started_at", "ended_at").
		From("dialogs").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("id DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get user dialogs: build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("get user dialogs: query: %w", err)
	}
	defer rows.Close()

	var res []models.Dialog
	for rows.Next() {
		var d models.Dialog
		if err := rows.Scan(
			&d.ID,
			&d.CaseID,
			&d.UserID,
			&d.ModelName,
			&d.StartedAt,
			&d.EndedAt,
		); err != nil {
			return nil, fmt.Errorf("get user dialogs: scan: %w", err)
		}
		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get user dialogs: rows iteration: %w", err)
	}

	return res, nil
}

func (r *PostgresDialogRepo) GetDialogsByCaseID(ctx context.Context, caseID int64, limit, offset int) ([]models.Dialog, error) {
	query := psql.
		Select("id", "case_id", "user_id", "model_name", "started_at", "ended_at").
		From("dialogs").
		Where(sq.Eq{"case_id": caseID}).
		OrderBy("id DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get dialogs by case id: build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("get dialogs by case id: query: %w", err)
	}
	defer rows.Close()

	var res []models.Dialog
	for rows.Next() {
		var d models.Dialog
		if err := rows.Scan(
			&d.ID,
			&d.CaseID,
			&d.UserID,
			&d.ModelName,
			&d.StartedAt,
			&d.EndedAt,
		); err != nil {
			return nil, fmt.Errorf("get dialogs by case id: scan: %w", err)
		}
		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get dialogs by case id: rows iteration: %w", err)
	}

	return res, nil
}

func (r *PostgresDialogRepo) CountDialogs(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM dialogs").Scan(&count)
	return count, err
}
