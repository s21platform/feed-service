//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"

	feed "github.com/s21platform/feed-proto/feed-proto"
)

type DBRepo interface {
	CreateUserPost(ctx context.Context, UUID string, in *feed.CreateUserPostIn) (*feed.CreateUserPostOut, error)
}
