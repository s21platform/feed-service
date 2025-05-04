package user

import "context"

type DBRepo interface {
	SaveNewEntity(ctx context.Context, UUID, metadata string) error
}
