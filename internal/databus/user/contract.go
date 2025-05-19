package user

import (
	"context"

	"github.com/s21platform/user-service/pkg/user"
)

type DBRepo interface {
	SaveNewEntity(ctx context.Context, UUID, metadata string) error
	SaveNewEntitySuggestion(ctx context.Context, followers []*user.Peer) error
}
