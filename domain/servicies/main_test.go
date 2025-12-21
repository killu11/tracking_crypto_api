package servicies_test

import (
	"crypto_api/api_client/geckocoin"
	"crypto_api/config"
	"crypto_api/domain/servicies"
	"crypto_api/infrastructure/persistence/postgres"
	"crypto_api/infrastructure/persistence/redis"
	"log"
	"os"
	"testing"
)

var coinService *servicies.CoinService

func TestMain(m *testing.M) {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := postgres.NewPostgresConnection(conf.Postgres)
	if err != nil {
		log.Fatalln(err)
	}
	redisClient, err := redis.NewRedisConnection(conf.Redis)
	if err != nil {
		log.Fatalln(err)
	}

	coinRepo := postgres.NewCoinRepository(db)
	coinCache := redis.NewCoinCacheRepository(redisClient)
	gecko := geckocoin.NewGeckoClient(conf.Gecko)
	trackingRepo := postgres.NewTrackingRepository(db)
	coinService = servicies.NewCoinService(coinRepo, coinCache, gecko, trackingRepo)

	os.Exit(m.Run())
}
