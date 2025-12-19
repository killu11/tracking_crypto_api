package entities

import (
	"errors"
	"time"
)

var ErrInvalidAmount = errors.New("amount less zero")

type Coin struct {
	ID           string    `json:"-"`
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	Usd          float64   `json:"usd"`
	LastUpdateAt time.Time `json:"last_updated_at"`
}

func (c *Coin) UpdatePrice(amount float64) error {
	if amount < 0 {
		return ErrInvalidAmount
	}
	c.Usd = amount
	return nil
}
