package postgres_test

import (
	"crypto_api/config"
	"crypto_api/domain/repositories"
	"crypto_api/persistence/postgres"
	"log"
	"os"
	"testing"
)

var coinRepo repositories.CoinRepository

func TestMain(m *testing.M) {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	db, err := postgres.NewPostgresConnection(conf.Postgres)
	if err != nil {
		log.Fatalln(err)
	}
	coinRepo = postgres.NewCoinRepository(db)
	os.Exit(m.Run())
}
