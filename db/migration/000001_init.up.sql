CREATE TABLE IF NOT EXISTS users
(
    id            SERIAL PRIMARY KEY,
    username      TEXT NOT NULL,
    email         TEXT NOT NULL,
    hash_password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS controllers
(
    id         SERIAL PRIMARY KEY,
    hw_key     TEXT NOT NULL UNIQUE,
    is_used    BOOLEAN DEFAULT FALSE,
    is_used_by INT     DEFAULT NULL
);