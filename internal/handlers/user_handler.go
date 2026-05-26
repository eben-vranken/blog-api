package handlers

import (
	"log"
	"net/http"

	"github.com/eben-vranken/blog-api/internal/repository"
)

type UserHandler struct {
	ur *repository.UserRepository
}

func (uh *UserHandler) Create(w http.ResponseWriter, req *http.Request) {
	log.Print("Creating")
}

func CreateUserHandler(ur *repository.UserRepository) UserHandler {
	t := new(UserHandler)
	t.ur = ur
	return *t
}
