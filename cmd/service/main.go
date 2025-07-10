package main

import (
	"fmt"
	"log"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"

	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/infra"
	db "github.com/s21platform/feed-service/internal/repository/postgres"
	"github.com/s21platform/feed-service/internal/service"
	"github.com/s21platform/feed-service/pkg/feed"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := db.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Fatalf("failed to create metrics object: %v", err)
	}

	feedService := service.New(dbRepo)
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
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
