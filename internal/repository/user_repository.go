package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/eben-vranken/blog-api/internal/auth"
	"github.com/eben-vranken/blog-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func (ur UserRepository) Create(ctx context.Context, user models.User) (models.UserResponse, error) {
	var userResponse models.UserResponse

	err := ur.db.QueryRowContext(ctx, `INSERT INTO users 
	(first_name,
	last_name,
	email,
	username,
	password_hash) 
	VALUES ($1, $2, $3, $4, $5)
	RETURNING user_id, first_name, last_name, email, username, created_at`,
		user.FirstName, user.LastName, user.Email, user.Username, user.PasswordHash,
	).Scan(&userResponse.UserID, &userResponse.FirstName, &userResponse.LastName, &userResponse.Email, &userResponse.Username, &userResponse.CreatedAt)

	return userResponse, err
}

func (ur UserRepository) GetAll(ctx context.Context) ([]models.UserResponse, error) {
	rows, err := ur.db.QueryContext(ctx, `SELECT
	user_id,
	first_name,
	last_name,
	email,
	username,
	created_at
	FROM users;`)

	var users []models.UserResponse = []models.UserResponse{}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user models.UserResponse

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

func (ur UserRepository) GetSpecific(ctx context.Context, id string) (models.UserResponse, error) {
	var user models.UserResponse

	err := ur.db.QueryRowContext(ctx, `SELECT
	user_id,
	first_name,
	last_name,
	email,
	username,
	created_at
	FROM users 
	WHERE user_id = $1;`, id).Scan(&user.UserID, &user.FirstName, &user.LastName, &user.Email, &user.Username, &user.CreatedAt)

	return user, err
}

func (ur UserRepository) Delete(ctx context.Context, id string, hashedPassword string) (sql.Result, error) {
	var user models.User

	err := ur.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, id).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	log.Print(auth.CheckPassword(user.PasswordHash, hashedPassword))

	if auth.CheckPassword(hashedPassword, user.PasswordHash) {
		result, err := ur.db.ExecContext(ctx, `DELETE FROM users WHERE user_id = $1`, id)

		return result, err
	}

	return nil, errors.New("Passwords do not match.")
}

func CreateUserRepository(db *sql.DB) UserRepository {
	t := new(UserRepository)
	t.db = db
	return *t
}
