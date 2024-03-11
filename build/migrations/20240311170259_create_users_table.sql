-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
    name varchar(100) NOT NULL,
    nick varchar(50) NOT NULL UNIQUE,
    email varchar(100) NOT NULL UNIQUE,
    password varchar(60) NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS followers;
-- +goose StatementEnd
