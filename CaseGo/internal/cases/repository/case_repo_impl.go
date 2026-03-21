package repository

import "database/sql"

type PostgresCaseRepo struct {
	db *sql.DB
}

func NewPostgresCaseRepo(db *sql.DB) *PostgresCaseRepo {
	return &PostgresCaseRepo{db: db}
}
