package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service     Service
	Postgres    Postgres
	Kafka       Kafka
	Metrics     Metrics
	Platform    Platform
	UserService UserService
	Logger      Logger
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
	UserTopic string `env:"USER_POST_CREATED"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

type UserService struct {
	Host string `env:"USER_SERVICE_HOST"`
	Port string `env:"USER_SERVICE_PORT"`
}

type Logger struct {
	Port string `env:"LOGGER_PORT"`
	Host string `env:"LOGGER_HOST"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
