package repository

import "database/sql"

type UserRepository struct {
	db *sql.DB
}

func CreateUserRepository(db *sql.DB) UserRepository {
	t := new(UserRepository)
	t.db = db
	return *t
}
