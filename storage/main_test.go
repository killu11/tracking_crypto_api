package storage_test

import (
	"crypto_api/config"
	"log"
	"os"
	"testing"
)

var (
	cfg, err = config.NewConfig()
)

func TestMain(m *testing.M) {
	if err == nil {
		os.Exit(m.Run())
	}
	log.Fatalf("load config: %v", err)
}
