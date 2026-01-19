package scheduler

import (
	"context"
	"crypto_api/domain/repositories"
	"log"
	"sync"
	"time"
)

type CleaningScheduler struct {
	duration time.Duration
	coinRepo repositories.CoinRepository
	m        *sync.Mutex
}

func (cs *CleaningScheduler) SetDuration(t time.Duration) {
	cs.m.Lock()
	cs.duration = t
	cs.m.Unlock()
}

func (cs *CleaningScheduler) Start(ctx context.Context) chan<- time.Duration {
	configChan := make(chan time.Duration)
	cs.run(ctx, configChan)
	return configChan
}

func (cs *CleaningScheduler) run(ctx context.Context, configChan <-chan time.Duration) {
	ticker := time.NewTicker(cs.duration)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case dur := <-configChan:
				if dur != cs.duration {
					cs.SetDuration(dur)
					ticker.Reset(dur)
					log.Printf("Cleaning scheduler interval updated to %v", dur)
				}
			case <-ticker.C:
				if err := cs.cleanUnactiveCoins(ctx); err != nil {
					log.Printf("|Warning| cleaning scheduler: %v\n", err)
				}
			}
		}
	}()
}
func (cs *CleaningScheduler) cleanUnactiveCoins(ctx context.Context) error {
	return cs.coinRepo.DeleteUnactiveCoins(ctx)
}
