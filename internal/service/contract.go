//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"
	"github.com/s21platform/feed-service/pkg/feed"
)

type DBRepo interface {
	FindTargetSuggestions(ctx context.Context, in *feed.GetFeedIn) ([]string, error)
	FindEntityInfo(ctx context.Context, targetSuggestions []string) (map[string][]string, error)
}
