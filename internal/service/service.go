package service

import (
	"context"
	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/pkg/feed"
	logger_lib "github.com/s21platform/logger-lib"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	feed.UnimplementedFeedServiceServer
	dbR DBRepo
}

func New(dbR DBRepo) *Service {
	return &Service{dbR: dbR}
}

func (s *Service) GetPost(ctx context.Context, in *feed.GetFeedIn) (*feed.GetFeedOut, error) {
	logger := logger_lib.FromContext(ctx, config.KeyLogger)
	logger.AddFuncName("GetPost")
	userUUID, ok := ctx.Value(config.KeyUUID).(string)

	if !ok || userUUID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "failed to retrieve uuid")
	}

	targetSuggestions, err := s.dbR.FindTargetSuggestions(ctx, in)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find target suggestions")
	}

	entityInfo, err := s.dbR.FindEntityInfo(ctx, targetSuggestions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to find entity info")
	}

	for k, v := range entityInfo {
		if k == "user" {
			// GetPostByUUIDS
		}
	}
}
