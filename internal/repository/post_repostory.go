package repository

import (
	"context"
	"database/sql"

	"github.com/eben-vranken/blog-api/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func (pr PostRepository) CreateDraft(ctx context.Context, post models.Post) (models.Post, error) {
	err := pr.db.QueryRowContext(ctx, `INSERT INTO posts (
	user_id,
	title,
	content,
	status
	) VALUES ($1, $2, $3, $4)
	RETURNING post_id, created_at, updated_at;
	`, &post.UserID, &post.Title, &post.Content, "draft").Scan(&post.PostID, &post.CreatedAt, &post.UpdatedAt)

	return post, err
}

func (pr PostRepository) PublishDraft(ctx context.Context, postId string) (models.Post, error) {
	var post models.Post
	
	err := pr.db.QueryRowContext(ctx, `UPDATE posts SET
	status = 'published',
	published_at = NOW()
	WHERE post_id = $1
	RETURNING post_id, user_id, title, content, status, created_at, published_at, updated_at;`, postId).Scan(&post.PostID, &post.UserID, &post.Title, &post.Content, &post.Status, &post.CreatedAt, &post.PublishedAt, &post.UpdatedAt)

	return post, err
}

func CreatePostRepository(db *sql.DB) PostRepository {
	t := new(PostRepository)
	t.db = db
	return *t
}
