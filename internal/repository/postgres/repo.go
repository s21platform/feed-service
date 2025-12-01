package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/model"
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

func (r *Repository) GetPostsByUserSubscriptions(ctx context.Context, userUUID string, limit int) ([]model.Entity, error) {
	query, args, err := squirrel.Select("e.external_uuid", "e.metadata").
		From("entities e").
		Join("entities_suggestion es ON es.post_uuid = e.uuid").
		Where(squirrel.Eq{"es.target_uuid": userUUID}).
		Where(squirrel.Eq{"e.metadata": "user"}).
		OrderBy("e.created_at DESC").
		Limit(uint64(limit)).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, fmt.Errorf("failed to build select query: %v", err)
	}

	var posts []model.Entity
	err = r.connection.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get posts by user subscriptions: %v", err)
	}

	return posts, nil
}
