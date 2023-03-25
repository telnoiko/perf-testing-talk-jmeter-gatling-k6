\c tasks;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    name     TEXT NOT NULL UNIQUE,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    tokens   TEXT[]
);

CREATE INDEX users_email_idx ON users (email);

CREATE TABLE tasks
(
    id          SERIAL PRIMARY KEY,
    description TEXT    NOT NULL,
    completed   BOOLEAN NOT NULL,
    user_id     INT     NOT NULL REFERENCES users (id) ON UPDATE CASCADE
);
