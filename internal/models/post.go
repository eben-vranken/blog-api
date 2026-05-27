package models

import "time"

type Post struct {
	PostID      int        `json:"post_id"`
	UserID      *int       `json:"user_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type PostRequest struct {
	PostID      int        `json:"post_id"`
	UserID      *int       `json:"user_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Status      Status     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	PublishedAt *time.Time `json:"published_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Password    string     `json:"password"`
}

type ProtectedPostRequest struct {
	UserID   int    `json:"user_id"`
	Password string `json:"password"`
}
