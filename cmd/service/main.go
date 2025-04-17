package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"

	feed "github.com/s21platform/feed-proto/feed-proto"
	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/infra"
	db "github.com/s21platform/feed-service/internal/repository/postgres"
	"github.com/s21platform/feed-service/internal/service"
	"google.golang.org/grpc"
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

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		log.Fatalf("cannot listen port: %s; Error: %v", cfg.Service.Port, err)
	}
	if err = server.Serve(lis); err != nil {
		log.Fatalf("cannot start grpc, port: %s; Error: %v", cfg.Service.Port, err)
	}
}
