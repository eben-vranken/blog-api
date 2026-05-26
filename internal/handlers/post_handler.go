package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/eben-vranken/blog-api/internal/models"
	"github.com/eben-vranken/blog-api/internal/repository"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostHandler struct {
	pr *repository.PostRepository
}

func (ph *PostHandler) CreateDraft(w http.ResponseWriter, req *http.Request) {
	var post models.Post

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&post)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}

	post, err = ph.pr.CreateDraft(req.Context(), post)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(post)

	if err != nil {
		log.Print(err)
	}
}

func (ph *PostHandler) PublishDraft(w http.ResponseWriter, req *http.Request) {
	post, err := ph.pr.PublishDraft(req.Context(), req.PathValue("id"))

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
			w.Write([]byte("404 - Post not found"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(post)

	if err != nil {
		log.Print(err)
	}
}

func CreatePostHandler(pr repository.PostRepository) PostHandler {
	t := new(PostHandler)
	t.pr = &pr
	return *t
}
