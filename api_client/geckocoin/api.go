package geckocoin

import (
	"context"
	"crypto_api/api_client/geckocoin/dto"
	"crypto_api/domain/entities"
	"crypto_api/pkg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

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

func (gc *GeckoClient) GetCoinID(ctx context.Context, symbol string) (string, error) {
	symbol = pkg.NormalizeSymbol(symbol)

	coinID, err := gc.searchCoinID(ctx, symbol)
	if err != nil {
		return "", err
	}

	if coinID == "" {
		return "", ErrCoinNotFound
	}

	return coinID, nil
}

func (gc *GeckoClient) GetCoinByID(ctx context.Context, coinID string) (*entities.Coin, error) {
	response, err := gc.FetchEndpoint(ctx, fmt.Sprintf("coins/%s", coinID), nil)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 512))
		return nil, NewAPIError(response.Status, response.StatusCode, string(message))
	}

	var meta dto.CoinMeta
	if err = json.NewDecoder(response.Body).Decode(&meta); err != nil {
		return nil, err
	}
	if meta.ID == "" {
		return nil, ErrCoinNotFound
	}
	return meta.ToEntity(), nil
}

func (gc *GeckoClient) searchCoinID(ctx context.Context, symbol string) (string, error) {
	params := url.Values{}
	params.Set("query", symbol)
	response, err := gc.FetchEndpoint(ctx, "search", params)
	if err != nil {
		return "", fmt.Errorf("get coin ID: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 512))
		return "", NewAPIError(response.Status, response.StatusCode, string(message))
	}

	dec := json.NewDecoder(response.Body)
	for {
		token, err := dec.Token()
		if err != nil {
			return "", err
		}
		if key, ok := token.(string); ok && key == "coins" {
			token, err = dec.Token()
			if delim, ok := token.(json.Delim); ok && delim == '[' {
				break
			}
		}
	}

	if !dec.More() {
		return "", ErrCoinNotFound
	}

	type CoinID struct {
		ID string `json:"id"`
	}

	var coinMeta CoinID
	if err = dec.Decode(&coinMeta); err != nil {
		return "", fmt.Errorf("json decode coin id: %w", err)
	}

	return coinMeta.ID, nil
}

func (gc *GeckoClient) RefreshCoinPrice(ctx context.Context, coin *entities.Coin) error {
	params := url.Values{}
	params.Set("vs_currencies", "usd")
	params.Set("ids", coin.ID)

	response, err := gc.FetchEndpoint(ctx, "simple/price", params)
	if err != nil {
		return fmt.Errorf("refresh coin price: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode >= http.StatusBadRequest {
		message, _ := io.ReadAll(io.LimitReader(response.Body, 512))
		return NewAPIError(response.Status, response.StatusCode, string(message))
	}

	dec := json.NewDecoder(response.Body)
	if !dec.More() {
		return ErrCoinNotFound
	}
	for {
		token, err := dec.Token()
		if err != nil {
			return fmt.Errorf("failed read token from body: %w", err)
		}
		if key, ok := token.(string); ok && key == coin.ID {
			break
		}
	}
	if err = dec.Decode(&coin); err != nil {
		return fmt.Errorf("decode fresh price data: %w", err)
	}
	return nil
}
