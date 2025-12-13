package storage

import (
	"context"
	"crypto_api/config"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client   *redis.Client
	cacheTTL time.Duration
}

func NewCacheRepository(client *redis.Client, cacheTTL time.Duration) *RedisCache {
	return &RedisCache{client: client, cacheTTL: cacheTTL}
}

func (r *RedisCache) SetCryptoID(ctx context.Context, symbol, id string) error {
	if err := r.client.Set(ctx, symbol, id, r.cacheTTL).Err(); err != nil {
		return fmt.Errorf("caching coin id: %w", err)
	}
	return nil
}

func (r *RedisCache) GetCryptoID(ctx context.Context, symbol string) (string, bool, error) {
	coinID, err := r.client.Get(ctx, symbol).Result()

	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("redis get coin id: %w", err)
	}
	return coinID, true, nil
}

func NewRedisConnection(cfg *config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		Username:     cfg.User,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("ping redis: %w", err)
	}
	return client, nil
}
