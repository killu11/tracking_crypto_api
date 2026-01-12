package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type CoinCache struct {
	client *redis.Client
	logger *zap.SugaredLogger
}

func NewCoinCacheRepository(client *redis.Client) *CoinCache {
	return &CoinCache{
		client: client,
		logger: zap.NewExample().Sugar(),
	}
}

func (r *CoinCache) SetCryptoID(ctx context.Context, symbol, id string, ttl time.Duration) {
	if err := r.client.Set(ctx, symbol, id, ttl).Err(); err != nil {
		r.logger.Warnf("redis caching coin id: %v", err)
	}
}

func (r *CoinCache) GetCryptoID(ctx context.Context, symbol string) (string, bool) {
	coinID, err := r.client.Get(ctx, symbol).Result()

	if errors.Is(err, redis.Nil) {
		return "", false
	}
	if err != nil {
		r.logger.Warnf("redis getting coin id: %v", err)
		return "", false
	}
	return coinID, true
}

func (r *CoinCache) DropCryptoID(ctx context.Context, symbol string) {
	r.logger.Warnf("drop crypto id: %v", r.client.Del(ctx, symbol).Err())
}

func (r *CoinCache) SetNotFoundStatus(ctx context.Context, symbol string, ttl time.Duration) {
	if err := r.client.Set(
		ctx,
		fmt.Sprintf("coin:not_found:%s", symbol),
		"",
		ttl,
	).Err(); err != nil {
		r.logger.Warnf("redis caching 'not found' status: %v", err)
	}
}

func (r *CoinCache) IsNotFound(ctx context.Context, symbol string) bool {
	_, err := r.client.Get(ctx, fmt.Sprintf("coin:not_found:%s", symbol)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false
		}
		r.logger.Warnf("redis getting 'not found' status: %v", err)
	}
	return true
}
