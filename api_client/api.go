package api_client

import (
	"context"
	"crypto_api/domain/entities"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CoinMeta struct {
	ID         string `json:"id"`
	Symbol     string `json:"symbol"`
	Name       string `json:"name"`
	MarketData struct {
		ValutePrices map[string]float64 `json:"current_price"`
	} `json:"market_data"`
	LastUpdate time.Time `json:"last_updated"`
}

func (m CoinMeta) convertCoin(value string) entities.Coin {
	return entities.Coin{
		Id:           m.ID,
		Symbol:       m.Symbol,
		Name:         m.Name,
		Usd:          m.MarketData.ValutePrices[value],
		LastUpdateAt: m.LastUpdate,
	}
}

// Все для api POST запроса на отслежние
func (gc *GeckoClient) GetCoinByID(ctx context.Context, coinID string) (entities.Coin, error) {
	response, err := gc.FetchEndpoint(ctx, fmt.Sprintf("coins/%s", coinID), nil)
	if err != nil {
		return entities.Coin{}, err
	}
	defer response.Body.Close()

	var meta CoinMeta
	if err = json.NewDecoder(response.Body).Decode(&meta); err != nil {
		return entities.Coin{}, err
	}

	if meta.ID == "" {
		return entities.Coin{}, ErrCoinNotFound
	}
	return meta.convertCoin("usd"), nil
}

func (gc *GeckoClient) GetCoinID(ctx context.Context, symbol string) (string, error) {
	symbol = normalizeSymbol(symbol)
	coinID, found, err := gc.cache.GetCryptoID(ctx, symbol)
	if err != nil {
		return "", err
	}

	if found {
		return coinID, nil
	}

	coinID, err = gc.searchCoinID(ctx, symbol)
	if err != nil {
		return "", err
	}

	if coinID == "" {
		return "", ErrCoinNotFound
	}

	if err = gc.cache.SetCryptoID(ctx, symbol, coinID); err != nil {
		return coinID, err
	}

	return coinID, nil
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

func normalizeSymbol(symbol string) string {
	return strings.ToUpper(strings.TrimSpace(symbol))
}
