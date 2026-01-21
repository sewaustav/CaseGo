package profilerepo

import (
	"context"
	"database/sql"
	"fmt"
)

func (r *PostgresProfileRepo) Begin(ctx context.Context) (*sql.Tx, error) {
    if db, ok := r.db.(*sql.DB); ok {
        return db.BeginTx(ctx, nil)
    }
    return nil, fmt.Errorf("repository must be initialized with *sql.DB to start transactions")
}

// WithTx — создаем копию репо, но с транзакцией внутри
func (r *PostgresProfileRepo) WithTx(tx *sql.Tx) ProfileRepo {
    return &PostgresProfileRepo{
        db: tx,
    }
}