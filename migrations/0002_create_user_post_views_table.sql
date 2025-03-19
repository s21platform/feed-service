-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS user_post_views
(
    uuid       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    post_uuid  UUID NOT NULL,
    user_uuid  UUID NOT NULL,
    view_at         TIMESTAMP,
    created_at      TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_uuid) REFERENCES user_posts (uuid)
);

-- +goose Down
DROP TABLE IF EXISTS user_post_views;
DROP EXTENSION IF EXISTS pgcrypto;