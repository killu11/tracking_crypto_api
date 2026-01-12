package scheduler

import (
	"database/sql"
	"time"
)

type RefreshScheduler struct {
	Duration time.Duration
	db       *sql.DB
}
