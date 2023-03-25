\c tasks;

CREATE TABLE users
(
    id       SERIAL PRIMARY KEY,
    name     TEXT NOT NULL,
    email    TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    created  TIMESTAMP NOT NULL DEFAULT NOW(),
    tokens   TEXT[] NOT NULL DEFAULT '{}'
);

CREATE INDEX users_email_idx ON users (email);

CREATE TABLE tasks
(
    id          SERIAL PRIMARY KEY,
    description TEXT    NOT NULL,
    completed   BOOLEAN NOT NULL,
    created     TIMESTAMP NOT NULL DEFAULT NOW(),
    owner       INT     NOT NULL REFERENCES users (id) ON UPDATE CASCADE
);
