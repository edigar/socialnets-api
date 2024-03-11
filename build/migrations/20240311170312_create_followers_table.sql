-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS followers (
    user_id uuid NOT NULL,
    follower uuid NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, follower)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS followers;
-- +goose StatementEnd
