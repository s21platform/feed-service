-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS user_posts
(
    uuid       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_uuid UUID NOT NULL,
    content         TEXT NOT NULL,
    created_at      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS user_posts;
DROP EXTENSION IF EXISTS pgcrypto;