package api_client_test

import (
	gecko_client "crypto_api/api_client"
	"crypto_api/config"
	"crypto_api/storage"
	"log"
	"os"
	"testing"
	"time"
)

var conf, confErr = config.NewConfig()
var gecko *gecko_client.GeckoClient

func TestMain(m *testing.M) {
	if confErr == nil {
		redisClient, err := storage.NewRedisConnection(conf.Redis)
		if err != nil {
			log.Fatalln(err)
		}
		
		gecko = gecko_client.NewGeckoClient(
			conf.Gecko,
			storage.NewCacheRepository(redisClient, 1*time.Hour),
		)
		os.Exit(m.Run())
	}
	log.Fatalln(confErr)
}
