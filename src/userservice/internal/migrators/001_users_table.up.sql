CREATE TABLE IF NOT EXISTS users (
    id serial not null,
    email TEXT UNIQUE,
    password_hash TEXT,
    created_at TIMESTAMP
);