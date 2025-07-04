-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS entities_suggestion
(
    uuid         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_uuid    UUID NOT NULL,
    target_uuid  UUID NOT NULL,
    view_at           TIMESTAMP,
    created_at        TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_uuid) REFERENCES entities (uuid)
);

-- +goose Down
DROP TABLE IF EXISTS entities_suggestion;
DROP EXTENSION IF EXISTS pgcrypto;