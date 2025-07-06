package user

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	userproto "github.com/s21platform/user-service/pkg/user"

	"github.com/s21platform/feed-service/internal/config"
)

type Service struct {
	client userproto.UserServiceClient
}

func NewService(cfg *config.Config) *Service {
	connStr := fmt.Sprintf("%s:%s", cfg.UserService.Host, cfg.UserService.Port)

	conn, err := grpc.NewClient(connStr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to user-service: %v", err)
	}

	client := userproto.NewUserServiceClient(conn)

	return &Service{client: client}
}

func (s *Service) GetWhoFollowPeer(ctx context.Context, userUUID string) ([]*userproto.Peer, error) {
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs("uuid", userUUID))

	resp, err := s.client.GetWhoFollowPeer(ctx, &userproto.GetWhoFollowPeerIn{Uuid: userUUID})
	if err != nil {
		return nil, fmt.Errorf("failed to get user followers: %v", err)
	}

	return resp.Subscribers, nil
}
