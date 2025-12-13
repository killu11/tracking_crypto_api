package geckocoin_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"
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
			value := reflect.ValueOf(coin)
			anyType := reflect.TypeOf(coin)
			for i := range anyType.NumField() {
				t.Logf("%s:%v", anyType.Field(i).Name, value.Field(i).Interface())
			}
		})
	}
}
