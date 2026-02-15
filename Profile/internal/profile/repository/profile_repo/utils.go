package profilerepo

import (
	"context"
	"errors"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/YoungFlores/Case_Go/Profile/internal/profile/models"
	repoerr "github.com/YoungFlores/Case_Go/Profile/internal/profile/repository/errors"
	"github.com/lib/pq"
)

func extractField(constraint string) string {
	constraint = strings.ToLower(constraint)
	if strings.Contains(constraint, "username") {
		return "username"
	}
	if strings.Contains(constraint, "user_id") {
		return "user_id"
	}
	return "unknown"
}

func (r *PostgresProfileRepo) fetchProfile(ctx context.Context, builder squirrel.SelectBuilder) (*models.Profile, error) {
	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var p models.Profile
	err = r.db.QueryRowContext(ctx, sql, args...).Scan(
		&p.ID, &p.UserID, &p.Avatar, &p.IsActive, &p.Description, &p.Username, &p.Name, &p.Surname,
		&p.Patronymic, &p.City, &p.Age, &p.Sex, &p.Profession, &p.CaseCount, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PostgresProfileRepo) handlePgError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return &repoerr.RepoError{
			Field: extractField(pgErr.Constraint),
			Err:   repoerr.ErrConflict,
		}
	}
	return err
}
