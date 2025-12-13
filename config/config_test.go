package config_test

import (
	"crypto_api/config"
	"testing"
)

func TestNewConfig(t *testing.T) {
	_, err := config.NewConfig()
	if err != nil {
		t.Error(err)
		return
	}
}
