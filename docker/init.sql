SELECT 'CREATE DATABASE socialnets'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'socialnets')\gexec
\c socialnets

DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    nick varchar(50) NOT NULL UNIQUE,
    email varchar(100) NOT NULL UNIQUE,
    password varchar(60) NOT NULL,
    created_at timestamp default current_timestamp,
    updated_at timestamp
);

CREATE TABLE followers (
    user_id int NOT NULL,
    follower int NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, follower)
);

CREATE TABLE posts (
    id serial PRIMARY KEY,
    title varchar(100) NOT NULL,
    content varchar(500) NOT NULL,
    author int NOT NULL,
    likes int DEFAULT 0,
    created_at timestamp default current_timestamp,
    FOREIGN KEY (author) REFERENCES users(id) ON DELETE CASCADE
);


