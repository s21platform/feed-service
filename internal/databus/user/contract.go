package user

import "context"

type DBRepo interface {
	CreateUserPost(ctx context.Context, postUUID string) error
}
