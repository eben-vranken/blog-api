-- +goose Up
CREATE TABLE post_tags (
    post_id INT NOT NULL REFERENCES posts(post_id) ON DELETE CASCADE,
    tag_id INT NOT NULL REFERENCES tags(tag_id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, tag_id)
);

-- +goose Down
DROP TABLE post_tags;