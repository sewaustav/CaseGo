package repository

import "database/sql"

type PostgresInteractionRepo struct {
	db DBTX
}

func NewPostgresInteractionRepo(db *sql.DB) *PostgresInteractionRepo {
	return &PostgresInteractionRepo{db: db}
}
