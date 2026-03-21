package repository

import "database/sql"

type PostgresDialogRepo struct {
	db *sql.DB
}

func NewPostgresDialogRepo(db *sql.DB) *PostgresDialogRepo {
	return &PostgresDialogRepo{db: db}
}
