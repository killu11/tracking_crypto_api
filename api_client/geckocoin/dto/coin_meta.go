package dto

import (
	"crypto_api/domain/entities"
	"crypto_api/pkg"
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

func (m CoinMeta) ToEntity() *entities.Coin {
	return &entities.Coin{
		ID:           m.ID,
		Symbol:       pkg.NormalizeSymbol(m.Symbol),
		Name:         m.Name,
		Usd:          m.MarketData.ValutePrices["usd"],
		LastUpdateAt: m.LastUpdate,
	}
}
