package dto

import "crypto_api/domain/entities"

type GetTrackableCoinResponse entities.Coin
type GetTrackableCoinsResponse struct {
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
