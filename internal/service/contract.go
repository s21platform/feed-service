//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	userproto "github.com/s21platform/user-service/pkg/user"

	"github.com/s21platform/feed-service/internal/model"
)

type DBRepo interface {
	GetPostsByUserSubscriptions(ctx context.Context, userUUID string, limit int) ([]model.Entity, error)
}

type UserClient interface {
	GetPeerFollow(ctx context.Context, userUUID string) ([]*userproto.Peer, error)
	GetPostsByIds(ctx context.Context, userUUID string, postUUIDs []string) ([]*userproto.PostInfo, error)
}
