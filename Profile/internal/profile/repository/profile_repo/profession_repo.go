package profilerepo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
)

func (r *PostgresProfileRepo) AddProfession(ctx context.Context, profession *models.UserProfession) (*models.UserProfession, error) {
	query := psql.Insert("users_categories").
		Columns("user_id", "profession_id").
		Values(profession.UserID, profession.ProfessionID).
		Suffix("RETURNING id, user_id, profession_id")
	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	if err := r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&profession.ID, &profession.UserID, &profession.ProfessionID); err != nil {
		return nil, r.handlePgError(err)
	}

	return profession, nil
}

func (r *PostgresProfileRepo) EditProfession(ctx context.Context, profession *models.UserProfession) (*models.UserProfession, error) {
	query := psql.Update("users_categories").
		Set("profession_id", profession.ProfessionID).
		Where(sq.Eq{"id": profession.ID}).
		Suffix("RETURNING id, user_id, profession_id")

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, sqlStr, args...).Scan(&profession.ID, &profession.UserID, &profession.ProfessionID)
	if err != nil {
		return nil, r.handlePgError(err)
	}

	return profession, nil

}

func (r *PostgresProfileRepo) DeleteProfession(ctx context.Context, professionID int64) error {
	query := psql.Delete("users_categories").Where(sq.Eq{"id": professionID})

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return err
	}

	res, err := r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	count, _ := res.RowsAffected()
	if count == 0 {
		return repoerr.ErrNotFound
	}

	return nil
}

func (r *PostgresProfileRepo) GetAllProfessions(ctx context.Context, userID int64) ([]models.UserProfession, error) {
	sqlStr, args, err := psql.Select("id", "user_id", "profession_id").From("users_categories").Where(sq.Eq{"user_id": userID}).ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var professions []models.UserProfession
	for rows.Next() {
		var profession models.UserProfession
		if err := rows.Scan(&profession.ID, &profession.UserID, &profession.ProfessionID); err != nil {
		}
	}
	return professions, nil
}

func (r *PostgresProfileRepo) GetProfileIDByProfessionID(ctx context.Context, professionID int64) (int64, error) {
	strSql, args, err := psql.Select("user_id").From("users_categories").Where(sq.Eq{"profession_id": professionID}).ToSql()
	if err != nil {
		return 0, err
	}

	var userID int64
	if err = r.db.QueryRowContext(ctx, strSql, args...).Scan(&userID); err != nil {
		return 0, err
	}
	return userID, nil

}
