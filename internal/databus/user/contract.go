package user

import (
	"context"

	"github.com/s21platform/user-service/pkg/user"
)

type DBRepo interface {
	SaveNewEntity(ctx context.Context, UUID, metadata string) (string, error)
	SaveNewEntitySuggestion(ctx context.Context, postUUID, followerUUID string) error
}

type UserClient interface {
	GetWhoFollowPeer(ctx context.Context, userUUID string) ([]*user.Peer, error)
}
