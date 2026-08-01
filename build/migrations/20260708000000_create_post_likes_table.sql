-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS post_likes (
    post_id    int  NOT NULL,
    user_id    uuid NOT NULL,
    created_at timestamp DEFAULT current_timestamp,
    PRIMARY KEY (post_id, user_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- +goose StatementEnd
-- +goose StatementBegin
ALTER TABLE posts DROP COLUMN IF EXISTS likes;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE posts ADD COLUMN IF NOT EXISTS likes int DEFAULT 0;
-- +goose StatementEnd
-- +goose StatementBegin
DROP TABLE IF EXISTS post_likes;
-- +goose StatementEnd
