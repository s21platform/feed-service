package main

import (
	"context"
	"log"

	_ "github.com/lib/pq"

	kafkalib "github.com/s21platform/kafka-lib"
	logger_lib "github.com/s21platform/logger-lib"
	"github.com/s21platform/metrics-lib/pkg"

	client "github.com/s21platform/feed-service/internal/client/user"
	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/databus/user"
	"github.com/s21platform/feed-service/internal/repository/postgres"
)

const userPostConsumerGroupID = "new-post-creater"

func main() {
	cfg := config.MustLoad()
	logger := logger_lib.New(cfg.Logger.Host, cfg.Logger.Port, cfg.Service.Name, cfg.Platform.Env)

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Println("failed to connect graphite: ", err)
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)
	ctx = context.WithValue(ctx, config.KeyLogger, logger)

	userConsumerConfig := kafkalib.DefaultConsumerConfig(
		cfg.Kafka.Host,
		cfg.Kafka.Port,
		cfg.Kafka.UserTopic,
		userPostConsumerGroupID,
	)

	userConsumer, err := kafkalib.NewConsumer(userConsumerConfig, metrics)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	userClient := client.NewService(cfg)

	userHandler := user.New(dbRepo, userClient)
	userConsumer.RegisterHandler(ctx, userHandler.Handler)

	<-ctx.Done()
}
