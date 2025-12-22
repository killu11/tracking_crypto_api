package repositories

import (
	"context"
	"crypto_api/domain/entities"
	"errors"
)

var (
	ErrCoinAlreadyTracking = errors.New("coin_already_tracking")
	ErrCoinNotTracking     = errors.New("coin_not_tracking")
	ErrZeroTrackableCoins  = errors.New("zero_trackable_coins")
)

type TrackingRepository interface {
	Add(ctx context.Context, userID int, symbol string) error
	Exists(ctx context.Context, userID int, symbol string) (bool, error)
	FindBySymbol(ctx context.Context, userID int, symbol string) (*entities.Coin, error)
	GetAll(ctx context.Context, userID int) ([]*entities.Coin, error)
	UpdatePrice(ctx context.Context, coin *entities.Coin, userID int) error
	GetPriceHistory(ctx context.Context, userID int, symbol string) ([]*entities.Price, error)
	GetStatsBySymbol(ctx context.Context, userID int, symbol string) (*entities.PriceStatistic, float64, error)
	Delete(ctx context.Context, userID int, symbol string) error
}
