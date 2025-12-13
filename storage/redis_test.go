package storage_test

import (
	"crypto_api/storage"
	"testing"
)

func TestNewRedisConnection(t *testing.T) {
	t.Logf("cfg: %v:", cfg.Redis)
	c, redisErr := storage.NewRedisConnection(cfg.Redis)
	if redisErr != nil {
		t.Error(redisErr)
		return
	}
	t.Logf("redis_client: %s", c.String())
}
