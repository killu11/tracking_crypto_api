package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type PostgresConfig struct {
	DSN         string
	MaxOpenCons int
}

type RedisConfig struct {
	Addr        string        `yaml:"Addr"`
	Password    string        `yaml:"password"`
	User        string        `yaml:"user"`
	DB          int           `yaml:"DB"`
	MaxRetries  int           `yaml:"MaxRetries"`
	DialTimeout time.Duration `yaml:"DialTimeout"`
	Timeout     time.Duration `yaml:"Timeout"`
}

func NewRedisConfig() (*RedisConfig, error) {
	f, err := os.Open(configPath + "redis_cfg.yml") // Путь для тестов geckocoin клиента
	if err != nil {
		return nil, fmt.Errorf("create redis config: %w", err)
	}

	rc := &RedisConfig{}
	if err = yaml.NewDecoder(f).Decode(&rc); err != nil {
		return nil, err
	}

	rc.Timeout *= time.Second
	rc.DialTimeout *= time.Second
	return rc, nil
}

func NewPostgresConfig() *PostgresConfig {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("PG_USER"),
		os.Getenv("PG_PASSWORD"),
		os.Getenv("PG_HOST"),
		os.Getenv("PG_PORT"),
		os.Getenv("PG_DB"),
		os.Getenv("PG_SSL"),
	)

	return &PostgresConfig{DSN: connStr, MaxOpenCons: 2}
}
