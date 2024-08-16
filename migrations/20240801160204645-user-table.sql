/*
 * Created at: 2024-08-01 16:02:04
 * Description: user table
 */


CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    password_hash bytea NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
