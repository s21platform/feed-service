-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS entities
(
    uuid          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_uuid UUID NOT NULL,
    metadata           TEXT NOT NULL,
    created_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS user_posts;
DROP EXTENSION IF EXISTS pgcrypto;