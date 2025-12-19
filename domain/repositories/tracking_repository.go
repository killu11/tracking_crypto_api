package repositories

import (
	"context"
	"errors"
)

var (
	ErrAlreadyTracking = errors.New("coin_already_tracking")
)

type TrackingRepository interface {
	Add(ctx context.Context, userID int, symbol string) error
	Exists(ctx context.Context, userID int, symbol string) (bool, error)
	Delete(ctx context.Context, userID int, symbol string) error
}
