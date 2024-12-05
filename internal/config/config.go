package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

// Config содержит конфигурацию приложения.
type Config struct {
	DBDSN       string `env:"DB_DSN,required"`
	GRPCPort    string `env:"GRPC_PORT" envDefault:"50051"`
	MetricsPort string `env:"METRICS_PORT" envDefault:"9091"`
}

// LoadConfig загружает конфигурацию из переменных окружения или файла .env.
func LoadConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			return nil, fmt.Errorf("failed to load .env file: %w", err)
		}
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return cfg, nil
}
