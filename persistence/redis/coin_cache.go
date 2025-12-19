package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

/*
TODO: настроить внутренний логгер
*/

type CoinCache struct {
	client *redis.Client
}

func NewCoinCacheRepository(client *redis.Client) *CoinCache {
	return &CoinCache{client: client}
}

func (r *CoinCache) SetCryptoID(ctx context.Context, symbol, id string, ttl time.Duration) error {
	if err := r.client.Set(ctx, symbol, id, ttl).Err(); err != nil {
		return fmt.Errorf("caching coin id: %w", err)
	}
	return nil
}

func (r *CoinCache) GetCryptoID(ctx context.Context, symbol string) (string, bool, error) {
	coinID, err := r.client.Get(ctx, symbol).Result()

	if errors.Is(err, redis.Nil) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("redis get coin id: %w", err)
	}
	return coinID, true, nil
}

func (r *CoinCache) SetNotFoundStatus(ctx context.Context, symbol string) error {
	if err := r.client.Set(
		ctx,
		fmt.Sprintf("coin:not_found:%s", symbol),
		"",
		10*time.Minute,
	).Err(); err != nil {
		return fmt.Errorf("caching coin:not_found: %w", err)
	}
	return nil
}

func (r *CoinCache) IsNotFound(ctx context.Context, symbol string) (bool, error) {
	_, err := r.client.Get(ctx, fmt.Sprintf("coin:not_found:%s", symbol)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return true, nil
		}
		return false, fmt.Errorf("redis get: %w", err)
	}
	return false, nil
}
