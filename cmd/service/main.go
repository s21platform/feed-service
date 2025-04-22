package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	feed "github.com/s21platform/feed-proto/feed-proto"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/infra"
	"github.com/s21platform/feed-service/internal/service"
	db "github.com/s21platform/feed-service/internal/repository/postgres"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := db.New(cfg)
	defer dbRepo.Close()

	feedService := service.New(dbRepo)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
		),
	)

	feed.RegisterFeedServiceServer(server, feedService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("failed to start TCP listener: %v", err)
	}
	if err = server.Serve(listener); err != nil {
		log.Fatalf("failed to start gRPC listener: %v", err)
	}
}
