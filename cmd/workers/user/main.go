package main

import (
	"context"
	"log"

	"github.com/s21platform/feed-service/internal/config"
	"github.com/s21platform/feed-service/internal/databus/user"
	"github.com/s21platform/feed-service/internal/repository/postgres"
	kafkalib "github.com/s21platform/kafka-lib"
	"github.com/s21platform/metrics-lib/pkg"
)

func main() {
	cfg := config.MustLoad()

	dbRepo := postgres.New(cfg)
	defer dbRepo.Close()

	metrics, err := pkg.NewMetrics(cfg.Metrics.Host, cfg.Metrics.Port, cfg.Service.Name, cfg.Platform.Env)
	if err != nil {
		log.Println("failed to connect graphite: ", err)
	}

	ctx := context.WithValue(context.Background(), config.KeyMetrics, metrics)

	userConsumerConfig := kafkalib.DefaultConsumerConfig(
		cfg.Kafka.Host,
		cfg.Kafka.Port,
		cfg.Kafka.UserTopic,
		"new-post-creater",
	)

	userConsumer, err := kafkalib.NewConsumer(userConsumerConfig, metrics)
	if err != nil {
		log.Fatalf("failed to create consumer: %v", err)
	}

	userHandler := user.New(dbRepo)
	userConsumer.RegisterHandler(ctx, userHandler.Handler)

	log.Println("Post consumer started successfully!")

	<-ctx.Done()
}
