package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	JWT        JWTConfig        `yaml:"jwt"`
	HttpServer HttpServerConfig `yaml:"http"`
	Postgres   PostgresConfig   `yaml:"postgres"`
	Kafka      KafkaConfig      `yaml:"kafka"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// load .env if exists
	if err := godotenv.Load(); err != nil {
		slog.Error(err.Error())
	}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}
