package geckocoin_test

import (
	"crypto_api/api_client/geckocoin"
	"crypto_api/config"
	"log"
	"os"
	"testing"
)

var (
	conf, confErr = config.NewConfig()
	gecko         *geckocoin.GeckoClient
)

func TestMain(m *testing.M) {
	if confErr == nil {
		gecko = geckocoin.NewGeckoClient(
			conf.Gecko,
		)
		os.Exit(m.Run())
	}
	log.Fatalln(confErr)
}
