package geckocoin_test

import (
	"context"
	"crypto_api/api_client/geckocoin"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestGeckoClient_GetCoinID(t *testing.T) {
	symbol := "eth"
	id, err := gecko.GetCoinID(context.Background(), symbol)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("coin ticker: %s:%s\n ", symbol, id)
}

func TestGeckoClient_GetCoinByID(t *testing.T) {
	coinsSymbols := []string{
		"btc",
		"eth",
		"doge",
	}
	ctx := context.Background()
	for _, symbol := range coinsSymbols {
		t.Run(fmt.Sprintf("test_%s", symbol), func(t *testing.T) {
			coinID, err := gecko.GetCoinID(ctx, symbol)
			if err != nil {
				t.Error(err)
				return
			}

			coin, err := gecko.GetCoinByID(ctx, coinID)
			if err != nil {
				t.Error(err)
				return
			}
			log.Println(coin)
			value := reflect.ValueOf(*coin)
			anyType := reflect.TypeOf(*coin)
			for i := range anyType.NumField() {
				t.Logf("%s:%v", anyType.Field(i).Name, value.Field(i).Interface())
			}
		})
	}
}

func TestGeckoClient_RefreshCoinPrice(t *testing.T) {
	coins := []string{"btc", "eth", "doge", "tramp", "luna", "froge"}
	for _, symbol := range coins {
		t.Run(fmt.Sprintf("update_price_%s", symbol), func(t *testing.T) {
			id, err := gecko.GetCoinID(context.Background(), symbol)
			if err != nil {
				if errors.Is(err, geckocoin.ErrCoinNotFound) {
					t.Errorf("invalid symbol: %v", err)
					return
				}
				t.Error(err)
				return
			}

			freshPrice, err := gecko.GetFreshCoinData(context.Background(), id)
			t.Log(time.Time(freshPrice.LastUpdated))
			if err != nil {
				t.Fatal(err)
				return
			}
			b, _ := json.Marshal(freshPrice)
			t.Logf("%s", b)
		})
	}
}
