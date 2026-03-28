package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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
		Suffix("RETURNING case_id, user_id, model_name, started_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&dialog.ID); err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		res = append(res, d)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func scanNullableTime(t *time.Time) any {
	return t
}

func (r *PostgresDialogRepo) _unused() error {
	return fmt.Errorf("unused")
}
