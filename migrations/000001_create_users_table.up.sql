CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    email varchar(255) NOT NULL UNIQUE,
    password bytea NOT NULL,
    created_at timestamp DEFAULT NOW()
);

CREATE TYPE duration AS ENUM ('weekly', 'monthly', '6 months', 'yearly');

CREATE TABLE IF NOT EXISTS subscriptions (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL,
    name varchar(255) NOT NULL,
    start_date timestamp NOT NULL,
    end_date timestamp NOT NULL,
    duration duration NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id)
);

