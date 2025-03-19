CREATE TABLE IF NOT EXISTS sessions (
    id uuid PRIMARY KEY,
    refresh_token varchar(512) NOT NULL,
    user_email varchar(255) NOT NULL,
    is_revoked bool DEFAULT false,
    created_at timestamp DEFAULT NOW(),
    expires_at timestamp NOT NULL
)
