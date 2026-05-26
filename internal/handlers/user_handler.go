package handlers

import (
	"encoding/json"
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
	var user models.User

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&user)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}

	password, err := auth.HashPassword(user.PasswordHash)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	user.PasswordHash = password

	user, err = uh.ur.Create(req.Context(), user)

	if err != nil {
		e := err.(*pgconn.PgError)

		if e.Code == "23505" {
			log.Print(err)
			w.WriteHeader(http.StatusConflict)
			w.Write([]byte("A user with this value already exists. Please check for duplicates."))
			return
		} else {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal server error"))
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)

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

func CreateUserHandler(ur *repository.UserRepository) UserHandler {
	t := new(UserHandler)
	t.ur = ur
	return *t
}
