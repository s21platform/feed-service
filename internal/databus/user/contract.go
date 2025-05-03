package user

import "context"

type DBRepo interface {
	CreatePost(ctx context.Context, postUUID, metadata string) error
}
