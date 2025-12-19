package postgres_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestCoinRepository_FindBySymbol(t *testing.T) {
	symbols := map[string]struct{}{
		"BTC":   {},
		"ETH":   {},
		"luna":  {},
		"tramp": {},
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Microsecond*5000)
	defer cancel()
	for symbol := range symbols {
		t.Run(fmt.Sprintf("find_%s", symbol), func(t *testing.T) {
			coin, err := coinRepo.FindBySymbol(ctx, symbol)
			if err != nil {
				t.Error(err)
				return
			}
			bytes, err := json.MarshalIndent(*coin, "", "\t")
			if err != nil {
				t.Error(err)
				return
			}
			t.Logf("%s\n", bytes)
		})
	}
}
