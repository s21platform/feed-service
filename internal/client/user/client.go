package client

import (
	"context"
	"fmt"
	"log"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/user-service/pkg/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	client user.UserServiceClient
}

func NewService(cfg *config.Config) *Service {
	connStr := fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port)

	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}

	client := user.NewUserServiceClient(conn)

	return &Service{client: client}
}

func (s *Service) GetWhoFollowPeer(ctx context.Context, userUUID string) ([]*user.Peer, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", userUUID))

	resp, err := s.client.GetWhoFollowPeer(ctx, &user.GetWhoFollowPeerIn{Uuid: userUUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user info from user-service: %v", err)
	}

	return resp.Subscribers, nil
}
