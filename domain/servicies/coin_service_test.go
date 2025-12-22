package servicies_test

import (
	"context"
	"crypto_api/domain/entities"
	"crypto_api/domain/repositories"
	"crypto_api/infrastructure/dto"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"
)

type testData struct {
	id     int
	symbol string
}

func TestCoinService_StopTrackCoin(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	err := coinService.StopTrackCoin(ctx, 1, "btc")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCoinService_TrackCoin(t *testing.T) {
	response, err := coinService.TrackCoin(context.Background(), "doge", 1)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := json.MarshalIndent(response, "", "\t")
	t.Logf("%s", b)
}

func TestCoinService_GetPriceHistory(t *testing.T) {
	data := []testData{
		{1, "btc"},
		{1, "eth"},
		{1, "unknown"},
	}
	for _, item := range data {
		t.Run(
			fmt.Sprintf("get_%s_history_by_user_%d", item.symbol, item.id),
			func(t *testing.T) {
				history, err := coinService.GetPriceHistory(context.Background(), item.id, item.symbol)
				if err != nil {
					if errors.Is(err, repositories.ErrCoinNotTracking) {
						t.Log(err)
						return
					}
					t.Fatal(err)
				}
				t.Log(history)
			},
		)
	}
}

func TestCoinService_GetCoinStats(t *testing.T) {
	want := dto.StatisticResponse{
		Symbol: "BTC",
		Price:  88584,
		PriceStatistic: &entities.PriceStatistic{
			Min:           87500,
			Max:           88584,
			Avg:           88042,
			Change:        1084,
			PercentChange: 1.22369728,
			Records:       2,
		},
	}
	got, err := coinService.GetCoinStats(context.Background(), 1, "BTC")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf(
		`Ожидаемо:
			min: %f, max: %f, avg: %f, change: %f, percent_change: %f`,
		want.Min,
		want.Max,
		want.Avg,
		want.Change,
		want.PercentChange,
	)
	t.Logf(
		`Получено:
			min: %f, max: %f, avg: %f, change: %f, percent_change: %f`,
		got.Min,
		got.Max,
		got.Avg,
		got.Change,
		got.PercentChange,
	)
}

func TestCoinService_RefreshTrackableCoin(t *testing.T) {
	coin, err := coinService.RefreshTrackableCoin(context.Background(), 1, "BTC")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(coin.Usd, coin.LastUpdateAt)

	_, err = coinService.RefreshTrackableCoin(context.Background(), 1, "dsafqewgegxds")
	t.Log(err)
}
