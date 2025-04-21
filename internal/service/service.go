package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	feed "github.com/s21platform/feed-proto/feed-proto"

	"github.com/s21platform/feed-service/internal/config"
)

type Service struct {
	feed.UnimplementedFeedServiceServer
	dbR DBRepo
}

func New(dbR DBRepo) *Service {
	return &Service{dbR: dbR}
}

func (s *Service) CreateUserPost(ctx context.Context, in *feed.CreateUserPostIn) (*feed.CreateUserPostOut, error) {
	ownerUUID, ok := ctx.Value(config.KeyUUID).(string)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "failed to retrieve uuid")
	}

	out, err := s.dbR.Post(ctx, ownerUUID, in.Content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create post: %v", err)
	}

	return &feed.CreateUserPostOut{PostUuid: out}, nil
}
