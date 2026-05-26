-- +goose Up
CREATE TABLE posts (
    post_id SERIAL PRIMARY KEY,
    user_id INT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    status status NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
        REFERENCES users(user_id)
        ON DELETE SET NULL
);

-- +goose Down
DROP TABLE posts;