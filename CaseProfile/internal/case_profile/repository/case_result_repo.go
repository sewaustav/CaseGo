package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sewaustav/CaseGoProfile/apperrors"
	"github.com/sewaustav/CaseGoProfile/internal/case_profile/models"
)

type CaseResultRepo interface {
	UpdateProfile(ctx context.Context, result *models.CaseProfile) error
	StoreProfile(ctx context.Context, result *models.CaseProfile) error
	AddResult(ctx context.Context, result *models.CaseResult) error
	GetResultByDialogID(ctx context.Context, dialogID int64) (*models.CaseResult, error)
	GetProfileByUserID(ctx context.Context, userID int64) (*models.CaseProfile, error)
	GetProfileByID(ctx context.Context, id int64) (*models.CaseProfile, error)
	GetHistoryBy(ctx context.Context, userID int64, from time.Time) ([]*models.CaseProfileHistory, error)
	DeleteResultByID(ctx context.Context, id int64) error
}

type PostgresCaseResultRepo struct {
	db *sql.DB
}

func NewPostgresCaseResultRepo(db *sql.DB) CaseResultRepo {
	return &PostgresCaseResultRepo{db: db}
}

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func (p *PostgresCaseResultRepo) UpdateProfile(ctx context.Context, result *models.CaseProfile) error {
	query, args, err := psql.
		Insert("case_profiles").
		Columns("user_id", "total_cases", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative").
		Values(result.UserID, result.TotalCases, result.Assertiveness, result.Empathy, result.ClarityCommunication, result.Resistance, result.Eloquence, result.Initiative).
		Suffix("ON CONFLICT (user_id) DO UPDATE SET total_cases = EXCLUDED.total_cases, assertiveness = EXCLUDED.assertiveness, empathy = EXCLUDED.empathy, clarity_communication = EXCLUDED.clarity_communication, resistance = EXCLUDED.resistance, eloquence = EXCLUDED.eloquence, initiative = EXCLUDED.initiative").
		ToSql()
	if err != nil {
		return fmt.Errorf("update profile: build query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("update profile: exec: %w", err)
	}

	return nil
}

func (p *PostgresCaseResultRepo) StoreProfile(ctx context.Context, result *models.CaseProfile) error {
	query, args, err := psql.
		Insert("case_profile_histories").
		Columns("user_id", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative", "actual_date").
		Values(result.UserID, result.Assertiveness, result.Empathy, result.ClarityCommunication, result.Resistance, result.Eloquence, result.Initiative, result.ChangedAt).
		ToSql()
	if err != nil {
		return fmt.Errorf("store profile history: build query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("store profile history: exec: %w", err)
	}

	return nil
}

func (p *PostgresCaseResultRepo) AddResult(ctx context.Context, result *models.CaseResult) error {
	query, args, err := psql.
		Insert("case_profile_results").
		Columns("user_id", "case_id", "dialog_id", "steps_count", "tokens_used", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative").
		Values(result.UserID, result.CaseID, result.DialogID, result.StepsCount, result.TokensUsed, result.Assertiveness, result.Empathy, result.ClarityCommunication, result.Resistance, result.Eloquence, result.Initiative).
		ToSql()
	if err != nil {
		return fmt.Errorf("add result: build query: %w", err)
	}

	_, err = p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("add result: exec: %w", err)
	}

	return nil
}

func (p *PostgresCaseResultRepo) GetResultByDialogID(ctx context.Context, dialogID int64) (*models.CaseResult, error) {
	query, args, err := psql.
		Select("id", "user_id", "case_id", "dialog_id", "steps_count", "tokens_used", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative").
		From("case_profile_results").
		Where(sq.Eq{"dialog_id": dialogID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("get result by dialog id: build query: %w", err)
	}

	result := &models.CaseResult{}
	err = p.db.QueryRowContext(ctx, query, args...).Scan(
		&result.ID, &result.UserID, &result.CaseID, &result.DialogID, &result.StepsCount, &result.TokensUsed,
		&result.Assertiveness, &result.Empathy, &result.ClarityCommunication, &result.Resistance, &result.Eloquence, &result.Initiative,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("result dialog_id=%d: %w", dialogID, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get result by dialog id: scan: %w", err)
	}

	return result, nil
}

func (p *PostgresCaseResultRepo) GetProfileByUserID(ctx context.Context, userID int64) (*models.CaseProfile, error) {
	query, args, err := psql.
		Select("id", "user_id", "total_cases", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative").
		From("case_profiles").
		Where(sq.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("get profile by user id: build query: %w", err)
	}

	profile := &models.CaseProfile{}
	err = p.db.QueryRowContext(ctx, query, args...).Scan(
		&profile.ID, &profile.UserID, &profile.TotalCases,
		&profile.Assertiveness, &profile.Empathy, &profile.ClarityCommunication,
		&profile.Resistance, &profile.Eloquence, &profile.Initiative,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile user_id=%d: %w", userID, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get profile by user id: scan: %w", err)
	}

	return profile, nil
}

func (p *PostgresCaseResultRepo) GetProfileByID(ctx context.Context, id int64) (*models.CaseProfile, error) {
	query, args, err := psql.
		Select("id", "user_id", "total_cases", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative").
		From("case_profiles").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("get profile by id: build query: %w", err)
	}

	profile := &models.CaseProfile{}
	err = p.db.QueryRowContext(ctx, query, args...).Scan(
		&profile.ID, &profile.UserID, &profile.TotalCases,
		&profile.Assertiveness, &profile.Empathy, &profile.ClarityCommunication,
		&profile.Resistance, &profile.Eloquence, &profile.Initiative,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("profile id=%d: %w", id, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get profile by id: scan: %w", err)
	}

	return profile, nil
}

func (p *PostgresCaseResultRepo) GetHistoryBy(ctx context.Context, userID int64, from time.Time) ([]*models.CaseProfileHistory, error) {
	query, args, err := psql.
		Select("id", "user_id", "assertiveness", "empathy", "clarity_communication", "resistance", "eloquence", "initiative", "actual_date").
		From("case_profile_histories").
		Where(sq.Eq{"user_id": userID}).
		Where(sq.GtOrEq{"actual_date": from}).
		OrderBy("actual_date DESC").
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("get history: build query: %w", err)
	}

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get history: query: %w", err)
	}
	defer rows.Close()

	var result []*models.CaseProfileHistory
	for rows.Next() {
		var profile models.CaseProfileHistory
		if err := rows.Scan(
			&profile.ID, &profile.UserID,
			&profile.Assertiveness, &profile.Empathy, &profile.ClarityCommunication,
			&profile.Resistance, &profile.Eloquence, &profile.Initiative, &profile.Date,
		); err != nil {
			return nil, fmt.Errorf("get history: scan: %w", err)
		}
		result = append(result, &profile)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get history: rows iteration: %w", err)
	}

	return result, nil
}

func (p *PostgresCaseResultRepo) DeleteResultByID(ctx context.Context, id int64) error {
	query, args, err := psql.
		Delete("case_profile_results").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("delete result: build query: %w", err)
	}

	res, err := p.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("delete result: exec: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("delete result: rows affected: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("result id=%d: %w", id, apperrors.ErrNotFound)
	}

	return nil
}
