package scheduler

import (
	"database/sql"
	"time"
)

type CleaningScheduler struct {
	Duration time.Duration
	db       *sql.DB
}
