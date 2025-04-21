//go:generate mockgen -destination=mock_contract_test.go -package=${GOPACKAGE} -source=contract.go
package service

import (
	"context"
)

type DBRepo interface {
	Post(ctx context.Context, uuid string, content string) (string, error)
}
