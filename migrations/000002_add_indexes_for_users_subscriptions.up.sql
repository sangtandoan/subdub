CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);
CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date ON subscriptions (end_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions (start_date);

-- add foregin key or add column -> alter table
-- add index -> create index
