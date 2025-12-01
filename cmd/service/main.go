package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soheilhy/cmux"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	"github.com/s21platform/feed-service/internal/api"
	"github.com/s21platform/feed-service/internal/client/user"
	"github.com/s21platform/feed-service/internal/config"
	generated "github.com/s21platform/feed-service/internal/generated"
	"github.com/s21platform/feed-service/internal/infra"
	db "github.com/s21platform/feed-service/internal/repository/postgres"
	"github.com/s21platform/feed-service/internal/service"
	"github.com/s21platform/feed-service/pkg/feed"
)

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)
	ctx := logger_lib.NewContext(context.Background(), logger)

	dbRepo := db.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("failed to create metrics object: %v", err))
		log.Fatalf("failed to create metrics object: %v", err)
	}

	userClient := user.NewService(cfg)
	feedService := service.New(dbRepo, userClient)

	grpcSrv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			infra.AuthInterceptor,
			infra.MetricsInterceptor(metrics),
		),
	)

	feed.RegisterFeedServiceServer(grpcSrv, feedService)

	restHandler := api.New(feedService)
	router := chi.NewRouter()
	router.Use(infra.AuthRequest)

	generated.HandlerFromMux(restHandler, router)
	httpServer := &http.Server{
		Handler: router,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Service.Port))
	if err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("cannot listen port; error: %s", err))
		log.Fatalf("Cannot listen port: %s; Error: %s", cfg.Service.Port, err)
	}

	m := cmux.New(lis)

	grpcListener := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpListener := m.Match(cmux.HTTP1Fast())

	g, _ := errgroup.WithContext(context.Background())

	logger_lib.Info(ctx, "starting server")

	g.Go(func() error {
		if err := grpcSrv.Serve(grpcListener); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
			return fmt.Errorf("gRPC server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := httpServer.Serve(httpListener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("HTTP server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		if err := m.Serve(); err != nil {
			return fmt.Errorf("cannot start service: %v", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		logger_lib.Error(ctx, fmt.Sprintf("server error: %v", err))
		log.Fatalf("Server error: %v", err)
	}
}
