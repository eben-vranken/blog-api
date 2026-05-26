package repository

import (
	"context"
	"database/sql"

	"github.com/eben-vranken/blog-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func (ur UserRepository) Create(ctx context.Context, user models.User) (models.User, error) {
	err := ur.db.QueryRowContext(ctx, `INSERT INTO users 
	(first_name,
	last_name,
	email,
	username,
	password_hash) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING user_id, created_at`,
		user.FirstName, user.LastName, user.Email, user.Username, user.PasswordHash,
	).Scan(&user.UserID, &user.CreatedAt)

	return user, err
}

func (ur UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	rows, err := ur.db.QueryContext(ctx, `SELECT
	user_id,
	first_name,
	last_name,
	email,
	username,
	created_at
	FROM users;`)

	var users []models.User

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.User

		err := rows.Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.CreatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, rows.Err()
}

func CreateUserRepository(db *sql.DB) UserRepository {
	t := new(UserRepository)
	t.db = db
	return *t
}
