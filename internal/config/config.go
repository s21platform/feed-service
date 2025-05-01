package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Kafka    Kafka
	Metrics  Metrics
	Platform Platform
}

type Service struct {
	Port string `env:"FEED_SERVICE_PORT"`
	Name string `env:"FEED_SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"FEED_SERVICE_POSTGRES_USER"`
	Password string `env:"FEED_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"FEED_SERVICE_POSTGRES_DB"`
	Host     string `env:"FEED_SERVICE_POSTGRES_HOST"`
	Port     string `env:"FEED_SERVICE_POSTGRES_PORT"`
}

type Kafka struct {
	Host      string `env:"KAFKA_HOST"`
	Port      string `env:"KAFKA_PORT"`
	UserTopic string `env:"USER_CREATE_NEW_POST"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
