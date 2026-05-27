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
	var post models.PostRequest

	decoder := json.NewDecoder(req.Body)

	err := decoder.Decode(&post)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad request"))
		return
	}

	createdPost, err := ph.pr.CreateDraft(req.Context(), post)

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
			w.Write([]byte("404 - Post not found"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdPost)

	if err != nil {
		log.Print(err)
	}
}

func (ph *PostHandler) PublishDraft(w http.ResponseWriter, req *http.Request) {
	var editRequest models.ProtectedPostRequest

	err := json.NewDecoder(req.Body).Decode(&editRequest)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	post, err := ph.pr.PublishDraft(req.Context(), req.PathValue("id"), editRequest)

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

func (ph *PostHandler) ArchivePost(w http.ResponseWriter, req *http.Request) {
	var editRequest models.ProtectedPostRequest

	err := json.NewDecoder(req.Body).Decode(&editRequest)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	post, err := ph.pr.ArchivePost(req.Context(), req.PathValue("id"), editRequest)

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

func (ph *PostHandler) GetAllPublished(w http.ResponseWriter, req *http.Request) {
	posts, err := ph.pr.GetAllPublished(req.Context())

	if err != nil {
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
	err = json.NewEncoder(w).Encode(posts)

	if err != nil {
		log.Print(err)
	}
}

func (ph *PostHandler) GetAllDrafts(w http.ResponseWriter, req *http.Request) {
	var postRequest models.ProtectedPostRequest

	err := json.NewDecoder(req.Body).Decode(&postRequest)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	posts, err := ph.pr.GetAllDrafts(req.Context(), postRequest)

	if err != nil {
		if errors.Is(err, repository.ErrInvalidPassword) {
			log.Print(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
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
	err = json.NewEncoder(w).Encode(posts)

	if err != nil {
		log.Print(err)
	}
}

func (ph *PostHandler) GetAllArchived(w http.ResponseWriter, req *http.Request) {
	var postRequest models.ProtectedPostRequest

	err := json.NewDecoder(req.Body).Decode(&postRequest)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	posts, err := ph.pr.GetAllArchived(req.Context(), postRequest)

	if err != nil {
		if errors.Is(err, repository.ErrInvalidPassword) {
			log.Print(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 - Unauthorized"))
			return
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
	err = json.NewEncoder(w).Encode(posts)

	if err != nil {
		log.Print(err)
	}
}

func (ph *PostHandler) DeletePost(w http.ResponseWriter, req *http.Request) {
	var editRequest models.ProtectedPostRequest

	err := json.NewDecoder(req.Body).Decode(&editRequest)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Bad Request"))
		return
	}

	result, err := ph.pr.DeletePost(req.Context(), req.PathValue("id"), editRequest)

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
			w.Write([]byte("404 - Post not found"))
			return
		}

		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	err = json.NewEncoder(w).Encode(result)

	if err != nil {
		log.Print(err)
	}
}

func CreatePostHandler(pr repository.PostRepository) PostHandler {
	t := new(PostHandler)
	t.pr = &pr
	return *t
}
