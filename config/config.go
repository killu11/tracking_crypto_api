package config

import (
	"fmt"

	"github.com/lpernett/godotenv"
)

type Config struct {
	Gecko    *GeckoApiConfig
	Postgres *PostgresConfig
	Redis    *RedisConfig
}

func NewConfig() (*Config, error) {
	if err := godotenv.Load("../cmd/.env"); err != nil {
		return nil, fmt.Errorf("failed create config: %w", err)
	}

	pgCfg := NewPostgresConfig()
	geckoCfg, err := NewGeckoApiConfig()
	if err != nil {
		return nil, err
	}
	redisCfg, err := NewRedisConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		Gecko:    geckoCfg,
		Postgres: pgCfg,
		Redis:    redisCfg,
	}, nil
}
