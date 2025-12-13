package api_client

import (
	"context"
	"crypto_api/config"
	"crypto_api/storage"
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
	cache  *storage.RedisCache
}

func NewGeckoClient(config *config.GeckoApiConfig, cacheRepo *storage.RedisCache) *GeckoClient {
	c := http.Client{
		Timeout: 2 * time.Second,
	}
	return &GeckoClient{
		Client: c,
		config: config,
		cache:  cacheRepo,
	}
}

func (gc *GeckoClient) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), gc.config.PingTimeout)
	defer cancel()

	response, err := gc.FetchEndpoint(ctx, "ping", nil)
	if err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("unexcepted ping status: %s", response.Status)
	}
	return nil
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
