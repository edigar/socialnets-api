-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS posts (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL,
    content varchar(500) NOT NULL,
    author uuid NOT NULL,
    likes int DEFAULT 0,
    created_at timestamp default current_timestamp,
    FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS posts;
-- +goose StatementEnd
