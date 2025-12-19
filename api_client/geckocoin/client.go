package geckocoin

import (
	"context"
	"crypto_api/config"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

/*
coins/market - возвращает список монет с рыночной капитализацией
https://docs.coingecko.com/v3.0.1/reference/coins-markets

coins/list - возвращает список поддерживаемых монет
https://docs.coingecko.com/v3.0.1/reference/coins-list
*/

type GeckoClient struct {
	http.Client
	config *config.GeckoApiConfig
}

func NewGeckoClient(config *config.GeckoApiConfig) *GeckoClient {
	c := http.Client{
		Timeout: 2 * time.Second,
	}
	return &GeckoClient{
		Client: c,
		config: config,
	}
}

// FetchEndpoint do get custom request to some geckoAPI endpoints
func (gc *GeckoClient) FetchEndpoint(ctx context.Context, path string, params url.Values) (*http.Response, error) {
	endpoint, err := gc.CreateEndpoint(path, params)
	if err != nil {
		return nil, err
	}
	r, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Add(gc.config.ApiHeader, gc.config.ApiKey)

	resp, err := gc.Do(r)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			return nil, err
		}
		return nil, fmt.Errorf("failed do request: %s: %w", r.URL.RawQuery, err)
	}
	return resp, nil
}

func (gc *GeckoClient) CreateEndpoint(suffix string, params url.Values) (string, error) {
	base, err := url.Parse(gc.config.BaseURL)
	if err != nil {
		return "", fmt.Errorf("parse url: %w", err)
	}
	endpoint, err := base.Parse(suffix)

	if err != nil {
		return "", fmt.Errorf("parse endpoint suffix: %w", err)
	}
	endpoint.RawQuery = params.Encode()
	return endpoint.String(), nil
}
