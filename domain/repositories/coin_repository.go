package repositories

import (
	"context"
	"crypto_api/domain/entities"
	"errors"
)

var (
	ErrCoinNotFound    = errors.New("coin_not_found")
	ErrUniqueViolation = errors.New("unique_constraint")
)

type CoinRepository interface {
	Save(ctx context.Context, coin *entities.Coin) error
	FindBySymbol(ctx context.Context, symbol string) (*entities.Coin, error)
	GetAll(ctx context.Context) ([]*entities.Coin, error)
	UpdatePrice(ctx context.Context, coin *entities.Coin) error
	Delete(ctx context.Context, symbol string) (bool, error)
}
