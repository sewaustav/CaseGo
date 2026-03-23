package repository

import (
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
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
		Columns("topic", "category", "is_generated", "description", "first_question", "creator", "created_at").
		Values(
			caseCopy.Topic,
			caseCopy.Category,
			caseCopy.IsGenerated,
			caseCopy.Description,
			caseCopy.FirstQuestion,
			caseCopy.Creator,
			caseCopy.CreatedAt,
		).
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&caseCopy.ID); err != nil {
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		res = append(res, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
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
		return nil, err
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&caseCopy.ID, &caseCopy.CreatedAt); err != nil {
		return nil, err
	}

	return caseCopy, nil
}

func (r *PostgresCaseRepo) DeleteCase(ctx context.Context, caseID int64) error {
	query := psql.Delete("cases").Where(sq.Eq{"id": caseID})
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("case not found")
	}
	return nil
}
