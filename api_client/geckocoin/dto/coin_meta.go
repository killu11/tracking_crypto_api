package dto

import (
	"crypto_api/domain/entities"
	"crypto_api/pkg"
	"encoding/json"
	"fmt"
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

type UnixTime time.Time

func (t *UnixTime) UnmarshalJSON(bytes []byte) error {
	var unixTimestamp int64
	if err := json.Unmarshal(bytes, &unixTimestamp); err != nil {
		return err
	}

	res := time.Unix(unixTimestamp, 0)
	if res.IsZero() {
		return fmt.Errorf("invalid unix time")
	}

	*t = UnixTime(res)
	return nil
}

func (t *UnixTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(*t))
}

type FreshPriceData struct {
	Symbol      string   `json:"-"`
	Price       float64  `json:"usd"`
	LastUpdated UnixTime `json:"last_updated_at"`
}
