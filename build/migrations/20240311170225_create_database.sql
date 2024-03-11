-- +goose Up
-- +goose StatementBegin
SELECT 'CREATE DATABASE socialnets'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'socialnets')\gexec
\c socialnets

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE IF EXISTS socialnets
-- +goose StatementEnd
