-- +goose Up
CREATE TYPE status AS ENUM ('draft', 'published', 'archived');

-- +goose Down
DROP TYPE status;