package repositories

import (
	"context"
	"crypto_api/domain/entities"
	"errors"
)

var (
	ErrCoinAlreadyTracking = errors.New("coin_already_tracking")
	ErrCoinNotTracking     = errors.New("coin_not_tracking")
)

type TrackingRepository interface {
	Add(ctx context.Context, userID int, symbol string) error
	Exists(ctx context.Context, userID int, symbol string) (bool, error)
	FindBySymbol(ctx context.Context, userID int, symbol string) (*entities.Coin, error)
	GetPriceHistory(ctx context.Context, userID int, symbol string) ([]*entities.Price, error)
	GetStatsBySymbol(ctx context.Context, userID int, symbol string) (*entities.PriceStatistic, float64, error)
	Delete(ctx context.Context, userID int, symbol string) error
}
