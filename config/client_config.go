package config

import (
	"fmt"
	"os"
	"time"
)

const (
	apiURL     = "https://api.coingecko.com/api/v3/"
	demoHeader = "x-cg-demo-api-key"
)

type GeckoApiConfig struct {
	BaseURL     string
	ApiKey      string
	ApiHeader   string
	PingTimeout time.Duration
}

func NewGeckoApiConfig() (*GeckoApiConfig, error) {
	ping := 1000 * time.Millisecond
	secret := os.Getenv("API_KEY")
	if secret == "" {
		return nil, fmt.Errorf("failed create gecko config: not found api-key")
	}
	return &GeckoApiConfig{
		BaseURL:     apiURL,
		ApiKey:      secret,
		ApiHeader:   demoHeader,
		PingTimeout: ping,
	}, nil
}
