package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/s21platform/feed-service/internal/config"
)

type Repository struct {
	connection *sqlx.DB
}

func New(cfg *config.Config) *Repository {
	conStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Database, cfg.Postgres.Host, cfg.Postgres.Port)

	conn, err := sqlx.Connect("postgres", conStr)
	if err != nil {
		log.Fatal("error connect: ", err)
	}

	return &Repository{
		connection: conn,
	}
}

func (r *Repository) Close() {
	_ = r.connection.Close()
}

func (r *Repository) SaveNewEntity(ctx context.Context, UUID, metadata string) (string, error) {
	query, args, err := squirrel.Insert("entities").
		Columns("external_uuid", "metadata").
		Values(UUID, metadata).
		Suffix("RETURNING uuid").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return "", fmt.Errorf("failed to build insert query: %v", err)
	}

	var postUUID string
	err = r.connection.GetContext(ctx, &postUUID, query, args...)

	if err != nil {
		return "", fmt.Errorf("failed to create user post in db: %v", err)
	}
	return postUUID, nil
}

func (r *Repository) SaveNewEntitySuggestion(ctx context.Context, postUUID, followerUUID string) error {
	query, args, err := squirrel.Insert("entities_suggestion").
		Columns("post_uuid", "target_uuid").
		Values(postUUID, followerUUID).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return fmt.Errorf("failed to build insert query: %v", err)
	}

	_, err = r.connection.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create user post in db: %v", err)
	}
	return nil
}
