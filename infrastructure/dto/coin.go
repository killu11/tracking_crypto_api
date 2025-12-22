package dto

import (
	"crypto_api/domain/entities"
)

type TrackableCoinResponse entities.Coin

func NewTrackableCoinResponse(coin *entities.Coin) *TrackableCoinResponse {
	return &TrackableCoinResponse{
		Symbol:       coin.Symbol,
		Name:         coin.Name,
		Usd:          coin.Usd,
		LastUpdateAt: coin.LastUpdateAt,
	}
}

type TrackableListResponse struct {
	Cryptos []*entities.Coin `json:"cryptos"`
}
type HistoryResponse struct {
	Symbol  string            `json:"symbol"`
	History []*entities.Price `json:"history"`
}

type StatisticResponse struct {
	Symbol                   string  `json:"symbol"`
	Price                    float64 `json:"current_price"`
	*entities.PriceStatistic `json:"stats"`
}
