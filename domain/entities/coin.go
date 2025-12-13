package entities

import (
	"time"
)

type UnixTime time.Time

type Coin struct {
	Id           string    `json:"id"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	Usd          float64   `json:"usd"`
	LastUpdateAt time.Time `json:"last_updated_at"`
}
