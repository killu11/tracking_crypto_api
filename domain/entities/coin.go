package entities

import (
	"errors"
	"time"
)

var (
	ErrInvalidAmount     = errors.New("amount less zero")
	ErrInvalidUpdateTime = errors.New("zero time err")
)

type Coin struct {
	ID           string    `json:"-"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	Usd          float64   `json:"current_price"`
	LastUpdateAt time.Time `json:"last_updated"`
}

func NewCoin(id, symbol string) *Coin {
	return &Coin{
		ID:     id,
		Symbol: symbol,
	}
}

type Price struct {
	Usd        float64   `json:"current_price"`
	LastUpdate time.Time `json:"last_updated"`
}

type PriceStatistic struct {
	Min           float64 `json:"min_price"`
	Max           float64 `json:"max_price"`
	Avg           float64 `json:"avg_price"`
	Change        float64 `json:"price_change,omitempty"`
	PercentChange float64 `json:"percent_price_change"`
	Records       int     `json:"records"`
}
