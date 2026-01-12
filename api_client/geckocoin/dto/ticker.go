package dto

import "time"

type UnixTime time.Time

type Ticker struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
}
