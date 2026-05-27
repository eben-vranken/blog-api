package repository

import (
	"context"
	"database/sql"

	"github.com/eben-vranken/blog-api/internal/auth"
	"github.com/eben-vranken/blog-api/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func (pr PostRepository) CreateDraft(ctx context.Context, post models.PostRequest) (*models.Post, error) {
	var user models.User

	err := pr.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, post.UserID).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(post.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	var createdPost models.Post

	err = pr.db.QueryRowContext(ctx, `INSERT INTO posts (
	user_id,
	title,
	content,
	status
	) VALUES ($1, $2, $3, $4)
	RETURNING post_id, user_id, title, content, status, created_at, updated_at;
	`, &post.UserID, &post.Title, &post.Content, "draft").Scan(&createdPost.PostID, &createdPost.UserID, &createdPost.Title, &createdPost.Content, &createdPost.Status, &createdPost.CreatedAt, &createdPost.UpdatedAt)

	return &createdPost, err
}

func (pr PostRepository) PublishDraft(ctx context.Context, postId string, editRequest models.ProtectedPostRequest) (*models.Post, error) {
	var user models.User

	err := pr.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, editRequest.UserID).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(editRequest.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	var post models.Post

	err = pr.db.QueryRowContext(ctx, `UPDATE posts SET
	status = 'published',
	published_at = NOW()
	WHERE post_id = $1 AND status = 'draft' AND user_id = $2
	RETURNING post_id, user_id, title, content, status, created_at, published_at, updated_at;`, postId, editRequest.UserID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.UpdatedAt)

	return &post, err
}

func (pr PostRepository) ArchivePost(ctx context.Context, postId string, editRequest models.ProtectedPostRequest) (*models.Post, error) {
	var user models.User

	err := pr.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, editRequest.UserID).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(editRequest.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	var post models.Post

	err = pr.db.QueryRowContext(ctx, `UPDATE posts SET
	status = 'archived',
	published_at = NOW()
	WHERE post_id = $1 AND status = 'published' AND user_id = $2
	RETURNING post_id, user_id, title, content, status, created_at, published_at, updated_at;`, postId, editRequest.UserID).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.UpdatedAt)

	return &post, err
}

func (pr PostRepository) GetAllPublished(ctx context.Context) ([]models.Post, error) {
	rows, err := pr.db.QueryContext(ctx, `SELECT 
	post_id,
	user_id,
	title,
	content,
	status,
	created_at,
	published_at,
	updated_at
	FROM posts
	WHERE status='published'`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post = []models.Post{}

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.UpdatedAt)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, err
}

func (pr PostRepository) GetAllDrafts(ctx context.Context, postRequest models.ProtectedPostRequest) ([]models.Post, error) {
	var user models.User

	err := pr.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, postRequest.UserID).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(postRequest.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	rows, err := pr.db.QueryContext(ctx, `SELECT 
	post_id,
	user_id,
	title,
	content,
	status,
	created_at,
	published_at,
	updated_at
	FROM posts
	WHERE status='draft'`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post = []models.Post{}

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.UpdatedAt)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, err
}

func (pr PostRepository) GetAllArchived(ctx context.Context, postRequest models.ProtectedPostRequest) ([]models.Post, error) {
	var user models.User

	err := pr.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, postRequest.UserID).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(postRequest.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	rows, err := pr.db.QueryContext(ctx, `SELECT 
	post_id,
	user_id,
	title,
	content,
	status,
	created_at,
	published_at,
	updated_at
	FROM posts
	WHERE status='archived'`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post = []models.Post{}

	for rows.Next() {
		var post models.Post

		err := rows.Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.UpdatedAt)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, err
}

func (pr PostRepository) DeletePost(ctx context.Context, postId string, editRequest models.ProtectedPostRequest) (sql.Result, error) {
	var user models.User

	err := pr.db.QueryRowContext(ctx, `SELECT
	password_hash
	FROM users 
	WHERE user_id = $1;`, editRequest.UserID).Scan(&user.PasswordHash)

	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(editRequest.Password, user.PasswordHash) {
		return nil, ErrInvalidPassword
	}

	result, err := pr.db.ExecContext(ctx, `DELETE FROM posts
	WHERE post_id = $1 AND user_id = $2
	`, postId, editRequest.UserID)

	return result, err
}

func CreatePostRepository(db *sql.DB) PostRepository {
	t := new(PostRepository)
	t.db = db
	return *t
}
