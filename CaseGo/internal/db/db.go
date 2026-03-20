package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type DataBase struct {
	db *sql.DB
}

func (d *DataBase) Open(dbName, dbUser, dbPassword, dbHost string, dbPort int) error {

	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost,
		dbPort,
		dbUser,
		dbPassword,
		dbName,
	)

	db, err := sql.Open(
		"postgres",
		psqlInfo,
	)

	if err != nil {
		db.Close()
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	d.db = db

	return nil
}

func (d *DataBase) Close() {
	d.db.Close()
}

func (d *DataBase) GetDB() *sql.DB {
	return d.db
}
