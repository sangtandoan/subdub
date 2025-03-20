CREATE TABLE IF NOT EXISTS auth_providers (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    provider varchar(255) NOT NULL,
    provider_id varchar(255) NOT NULL,
    created_at timestamp DEFAULT NOW(),

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);
