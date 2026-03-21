package repository

import "database/sql"

type PostgresInteractionRepo struct {
	db *sql.DB
}

func NewPostgresInteractionRepo(db *sql.DB) *PostgresInteractionRepo {
	return &PostgresInteractionRepo{db: db}
}
