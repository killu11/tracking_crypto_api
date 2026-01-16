package scheduler

import (
	"crypto_api/api_client/geckocoin"
	"database/sql"
	"time"
)

type RefreshPriceScheduler struct {
	duration    time.Duration
	db          *sql.DB
	geckoClient *geckocoin.GeckoClient
}
