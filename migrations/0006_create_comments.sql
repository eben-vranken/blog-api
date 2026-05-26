-- +goose Up
CREATE TABLE comments (
    comment_id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    user_id INT NULL,
    post_id INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(user_id)
        ON DELETE SET NULL,
    CONSTRAINT fk_post
        FOREIGN KEY(post_id)
        REFERENCES posts(post_id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE comments;