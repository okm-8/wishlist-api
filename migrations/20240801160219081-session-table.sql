/*
 * Created at: 2024-08-01 16:02:19
 * Description: session table
 */

CREATE TABLE IF NOT EXISTS sessions (
    id uuid PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    expire_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    FOREIGN KEY (id) REFERENCES sessions (id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS user_sessions_user_id_idx ON user_sessions (user_id);

CREATE TABLE IF NOT EXISTS signup_sessions (
    id uuid PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    FOREIGN KEY (id) REFERENCES sessions (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS signup_sessions_email_idx ON signup_sessions (email);

CREATE OR REPLACE VIEW user_sessions_view AS (
    SELECT
        session.id AS id,
        "user".id AS user_id,
        "user".email AS user_email,
        "user".name AS user_name,
        "user".admin AS user_admin,
        "user".password_hash AS user_password_hash,
        session.created_at AS created_at,
        session.expire_at AS expire_at
    FROM sessions AS session
    INNER JOIN user_sessions us on session.id = us.id
    INNER JOIN users AS "user" ON us.user_id = "user".id
);

CREATE OR REPLACE VIEW signup_sessions_view AS (
    SELECT
        session.id AS id,
        signup.email AS email,
        session.expire_at AS expire_at,
        session.created_at AS created_at
    FROM sessions AS session
    INNER JOIN signup_sessions signup on session.id = signup.id
);
