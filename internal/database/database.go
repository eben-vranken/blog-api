package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func New(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseUrl)

	if err != nil {
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, err
}
