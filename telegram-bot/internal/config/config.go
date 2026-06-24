package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Log      LogConfig
	Telegram TelegramConfig
	GRPC     GRCPConfig
}

func Load(path ...string) (*Config, error) {
	// Load only the .env files that actually exist. A missing file is not an
	// error: configuration may be provided entirely through the environment
	// (e.g. via docker-compose env_file), in which case no .env is present.
	var existing []string
	for _, p := range path {
		if p == "" {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			existing = append(existing, p)
		}
	}

	if len(existing) > 0 {
		if err := godotenv.Load(existing...); err != nil {
			return nil, err
		}
	}

	cfg := Config{}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
