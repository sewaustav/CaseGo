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

type PostgresInteractionRepo struct {
	db DBTX
	tx Tx
}

func NewPostgresInteractionRepo(db *sql.DB) *PostgresInteractionRepo {
	return &PostgresInteractionRepo{db: db}
}

func (r *PostgresInteractionRepo) current() DBTX {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *PostgresInteractionRepo) Begin(ctx context.Context) (Tx, error) {
	base, ok := r.db.(*sql.DB)
	if !ok {
		return nil, fmt.Errorf("begin: underlying db is not *sql.DB")
	}

	tx, err := base.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin tx: %w", err)
	}

	return tx, nil
}

func (r *PostgresInteractionRepo) WithTx(tx Tx) Interaction {
	return &PostgresInteractionRepo{
		db: r.db,
		tx: tx,
	}
}

func (r *PostgresInteractionRepo) PutInteraction(ctx context.Context, interaction *models.Interaction) error {
	query := psql.
		Insert("interactions").
		Columns("dialog_id", "step", "question", "answer", "tokens_used", "created_at").
		Values(interaction.DialogID, interaction.Step, interaction.Question, interaction.Answer, interaction.TokensUsed, interaction.CreatedAt).
		Suffix("RETURNING id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("put interaction: build query: %w", err)
	}

	if err := r.current().QueryRowContext(ctx, sqlStr, args...).Scan(&interaction.ID); err != nil {
		return fmt.Errorf("put interaction: exec: %w", err)
	}

	return nil
}

func (r *PostgresInteractionRepo) GetInteractionByID(ctx context.Context, interactionID int64) (*models.Interaction, error) {
	query := psql.
		Select("id", "dialog_id", "step", "question", "answer", "tokens_used", "created_at").
		From("interactions").
		Where(sq.Eq{"id": interactionID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get interaction by id: build query: %w", err)
	}

	var i models.Interaction
	if err := r.current().QueryRowContext(ctx, sqlStr, args...).Scan(
		&i.ID,
		&i.DialogID,
		&i.Step,
		&i.Question,
		&i.Answer,
		&i.TokensUsed,
		&i.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("interaction id=%d: %w", interactionID, apperrors.ErrNotFound)
		}
		return nil, fmt.Errorf("get interaction by id: scan: %w", err)
	}

	return &i, nil
}

func (r *PostgresInteractionRepo) GetInteractionsByDialogID(ctx context.Context, dialogID int64) ([]models.Interaction, error) {
	query := psql.
		Select("id", "dialog_id", "step", "question", "answer", "tokens_used", "created_at").
		From("interactions").
		Where(sq.Eq{"dialog_id": dialogID}).
		OrderBy("step ASC", "id ASC")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("get interactions by dialog id: build query: %w", err)
	}

	rows, err := r.current().QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, fmt.Errorf("get interactions by dialog id: query: %w", err)
	}
	defer rows.Close()

	var res []models.Interaction
	for rows.Next() {
		var i models.Interaction
		if err := rows.Scan(
			&i.ID,
			&i.DialogID,
			&i.Step,
			&i.Question,
			&i.Answer,
			&i.TokensUsed,
			&i.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("get interactions by dialog id: scan: %w", err)
		}
		res = append(res, i)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get interactions by dialog id: rows iteration: %w", err)
	}

	return res, nil
}
