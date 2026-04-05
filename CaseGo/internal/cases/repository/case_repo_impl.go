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

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

type PostgresCaseRepo struct {
	db *sql.DB
}

func NewPostgresCaseRepo(db *sql.DB) *PostgresCaseRepo {
	return &PostgresCaseRepo{db: db}
}

func (r *PostgresCaseRepo) CreateCase(ctx context.Context, caseCopy *models.Case) (*models.Case, error) {
	query := psql.
		Insert("cases").
		Columns("topic", "category", "is_generated", "description", "first_question", "creator").
		Values(
			caseCopy.Topic,
			caseCopy.Category,
			caseCopy.IsGenerated,
			caseCopy.Description,
			caseCopy.FirstQuestion,
			caseCopy.Creator,
		).
		Suffix("RETURNING id, created_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("create case: build query: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&caseCopy.ID, &caseCopy.CreatedAt); err != nil {
		return nil, fmt.Errorf("create case: exec: %w", err)
	}

	return caseCopy, nil
}

func (r *PostgresCaseRepo) GetCaseByID(ctx context.Context, caseID int64) (*models.Case, error) {
	query := psql.
		Select("id", "topic", "category", "is_generated", "description", "first_question", "creator", "created_at").
		From("cases").
		Where(sq.Eq{"id": caseID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get case by id: build query: %w", err)
	}

	var c models.Case
	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(
		&c.ID,
		&c.Topic,
		&c.Category,
		&c.IsGenerated,
		&c.Description,
		&c.FirstQuestion,
		&c.Creator,
		&c.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("case id=%d: %w", caseID, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get case by id: scan: %w", err)
	}

	return &c, nil
}

func (r *PostgresCaseRepo) GetCases(ctx context.Context, limit int, offset int) ([]models.Case, error) {
	query := psql.
		Select("id", "topic", "category", "is_generated", "description", "first_question", "creator", "created_at").
		From("cases").
		OrderBy("id DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get cases: build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("get cases: query: %w", err)
	}
	defer rows.Close()

	var res []models.Case
	for rows.Next() {
		var c models.Case
		if err := rows.Scan(
			&c.ID,
			&c.Topic,
			&c.Category,
			&c.IsGenerated,
			&c.Description,
			&c.FirstQuestion,
			&c.Creator,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get cases: scan: %w", err)
		}
		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get cases: rows iteration: %w", err)
	}

	return res, nil
}

func (r *PostgresCaseRepo) GetCasesByCategory(ctx context.Context, categoryID int32, limit int, offset int) ([]models.Case, error) {
	query := psql.
		Select("id", "topic", "category", "is_generated", "description", "first_question", "creator", "created_at").
		From("cases").
		Where(sq.Eq{"category": categoryID}).
		OrderBy("id DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get cases by category: build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("get cases by category: query: %w", err)
	}
	defer rows.Close()

	var res []models.Case
	for rows.Next() {
		var c models.Case
		if err := rows.Scan(
			&c.ID,
			&c.Topic,
			&c.Category,
			&c.IsGenerated,
			&c.Description,
			&c.FirstQuestion,
			&c.Creator,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get cases by category: scan: %w", err)
		}
		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get cases by category: rows iteration: %w", err)
	}

	return res, nil
}

func (r *PostgresCaseRepo) GetCasesByTopic(ctx context.Context, topicID string, limit int, offset int) ([]models.Case, error) {
	query := psql.
		Select("id", "topic", "category", "is_generated", "description", "first_question", "creator", "created_at").
		From("cases").
		Where(sq.Eq{"topic": topicID}).
		OrderBy("id DESC").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get cases by topic: build query: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("get cases by topic: query: %w", err)
	}
	defer rows.Close()

	var res []models.Case
	for rows.Next() {
		var c models.Case
		if err := rows.Scan(
			&c.ID,
			&c.Topic,
			&c.Category,
			&c.IsGenerated,
			&c.Description,
			&c.FirstQuestion,
			&c.Creator,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get cases by topic: scan: %w", err)
		}
		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get cases by topic: rows iteration: %w", err)
	}

	return res, nil
}

func (r *PostgresCaseRepo) PatchCase(ctx context.Context, caseCopy *models.Case) (*models.Case, error) {
	query := psql.
		Update("cases").
		Set("topic", caseCopy.Topic).
		Set("category", caseCopy.Category).
		Set("is_generated", caseCopy.IsGenerated).
		Set("description", caseCopy.Description).
		Set("first_question", caseCopy.FirstQuestion).
		Set("creator", caseCopy.Creator).
		Where(sq.Eq{"id": caseCopy.ID}).
		Suffix("RETURNING id, created_at")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("patch case: build query: %w", err)
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&caseCopy.ID, &caseCopy.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("case id=%d: %w", caseCopy.ID, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("patch case: scan: %w", err)
	}

	return caseCopy, nil
}

func (r *PostgresCaseRepo) DeleteCase(ctx context.Context, caseID int64) error {
	query := psql.Delete("cases").Where(sq.Eq{"id": caseID})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("delete case: build query: %w", err)
	}

	res, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("delete case: exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete case: rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("case id=%d: %w", caseID, apperrors.ErrNotFound)
	}

	return nil
}
