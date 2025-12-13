package geckocoin_test

import (
	"crypto_api/api_client/geckocoin"
	"crypto_api/config"
	"crypto_api/storage"
	"log"
	"os"
	"testing"
	"time"
)

var (
	conf, confErr = config.NewConfig()
	gecko         *geckocoin.GeckoClient
)

func TestMain(m *testing.M) {
	if confErr == nil {
		redisClient, err := storage.NewRedisConnection(conf.Redis)
		if err != nil {
			log.Fatalln(err)
		}

		gecko = geckocoin.NewGeckoClient(
			conf.Gecko,
			storage.NewCacheRepository(redisClient, 1*time.Hour),
		)
		os.Exit(m.Run())
	}
	log.Fatalln(confErr)
}
