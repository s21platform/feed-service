package postgres

import (
	"context"
	"fmt"
	"log"

	feed "github.com/s21platform/feed-proto/feed-proto"

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

func (r *Repository) CreateUserPost(ctx context.Context, UUID string, in *feed.CreateUserPostIn) (*feed.CreateUserPostOut, error) {
	var postObj model.Post

	postObj, err := postObj.PostToDTO(UUID, in)
	if err != nil {
		return nil, fmt.Errorf("failed toconvert grpc message to dto: %v", err)
	}

	query := squirrel.Insert("user_posts").
		Columns("owner_uuid", "content").
		Values(postObj.OwnerUUID, postObj.Content).
		Suffix("RETURNING uuid").
		PlaceholderFormat(squirrel.Dollar)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %v", err)
	}

	var newUUID string
	err = r.connection.QueryRowContext(ctx, sql, args...).Scan(&newUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to create post: %v", err)
	}

	return &feed.CreateUserPostOut{PostUuid: newUUID}, nil
}
