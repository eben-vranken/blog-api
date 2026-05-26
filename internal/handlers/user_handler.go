package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eben-vranken/blog-api/internal/auth"
	"github.com/eben-vranken/blog-api/internal/models"
	"github.com/eben-vranken/blog-api/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type UserHandler struct {
	ur *repository.UserRepository
}

func (uh *UserHandler) Create(w http.ResponseWriter, req *http.Request) {
	var user models.UserRequest

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&user)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}

	password, err := auth.HashPassword(user.Password)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var internalUser models.User

	internalUser.FirstName = user.FirstName
	internalUser.LastName = user.LastName
	internalUser.Email = user.Email
	internalUser.Username = user.Username
	internalUser.PasswordHash = password

	userResponse, err := uh.ur.Create(req.Context(), internalUser)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				log.Print(err)
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("A user with this value already exists. Please check for duplicates."))
				return
			}
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal server error"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(userResponse)

	if err != nil {
		log.Print(err)
		log.Print("500 - Internal server error")
	}
}

func (uh *UserHandler) GetAll(w http.ResponseWriter, req *http.Request) {
	users, err := uh.ur.GetAll(req.Context())

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(users)

	if err != nil {
		log.Print(err)
		log.Print("500 - Internal server error")
	}
}

func (uh *UserHandler) GetSpecific(w http.ResponseWriter, req *http.Request) {
	user, err := uh.ur.GetSpecific(req.Context(), req.PathValue("id"))

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			if pgErr.Code == "22P02" {
				log.Print(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Bad request\nID must be an integer"))
				return
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			log.Print(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - User not found"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)

	if err != nil {
		log.Print(err)
	}
}

func (uh *UserHandler) Delete(w http.ResponseWriter, req *http.Request) {
	var password models.UserPassword

	err := json.NewDecoder(req.Body).Decode(&password)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = uh.ur.Delete(req.Context(), req.PathValue("id"), password.Password)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.Is(err, repository.ErrInvalidPassword) {
			log.Print(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
		}

		if errors.As(err, &pgErr) {
			if pgErr.Code == "22P02" {
				log.Print(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Bad request\nID must be an integer"))
				return
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			log.Print(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - User not found"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) Update(w http.ResponseWriter, req *http.Request) {
	var userInfoToUpdate models.UserUpdateRequest

	err := json.NewDecoder(req.Body).Decode(&userInfoToUpdate)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = uh.ur.Update(req.Context(), req.PathValue("id"), userInfoToUpdate)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.Is(err, repository.ErrInvalidPassword) {
			log.Print(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
		}

		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				log.Print(err)
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte("A user with this value already exists. Please check for duplicates."))
				return
			}

			if pgErr.Code == "22P02" {
				log.Print(err)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("400 - Bad request\nID must be an integer"))
				return
			}
		}

		if errors.Is(err, sql.ErrNoRows) {
			log.Print(err)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - User not found"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func CreateUserHandler(ur *repository.UserRepository) UserHandler {
	t := new(UserHandler)
	t.ur = ur
	return *t
}
