package dto

import (
	"crypto_api/domain/entities"
	"encoding/json"
	"fmt"
	"time"
)

type FreshPriceData struct {
	Symbol      string   `json:"-"`
	Price       float64  `json:"usd"`
	LastUpdated UnixTime `json:"last_updated_at"`
}

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

func (d *FreshPriceData) ToEntity(symbol string) *entities.Coin {
	return &entities.Coin{
		Symbol:       symbol,
		Usd:          d.Price,
		LastUpdateAt: time.Time(d.LastUpdated),
	}
}
