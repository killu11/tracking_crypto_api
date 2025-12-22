package postgres

import (
	"context"
	"crypto_api/domain/entities"
	"crypto_api/domain/repositories"
	"database/sql"
	"errors"
	"fmt"
)

type CoinRepository struct {
	db *sql.DB
}

func (c *CoinRepository) Save(ctx context.Context, coin *entities.Coin) error {
	_, err := c.db.ExecContext(
		ctx,
		`INSERT INTO coins(symbol, name, current_price, last_updated) 
				VALUES ($1, $2, $3, $4) 
				ON CONFLICT DO NOTHING`,
		coin.Symbol, coin.Name, coin.Usd, coin.LastUpdateAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *CoinRepository) FindBySymbol(ctx context.Context, symbol string) (*entities.Coin, error) {
	var coin entities.Coin
	err := c.db.QueryRowContext(
		ctx,
		`SELECT symbol, name, current_price, last_updated
				FROM coins
				WHERE symbol = $1`,
		symbol,
	).Scan(&coin.Symbol, &coin.Name, &coin.Usd, &coin.LastUpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrCoinNotFound
		}
		return nil, fmt.Errorf("failed find coin by symbol: %w", err)
	}
	return &coin, nil
}

func (c *CoinRepository) GetAll(ctx context.Context) ([]*entities.Coin, error) {
	//TODO implement me
	panic("implement me")
}

func (c *CoinRepository) UpdatePrice(ctx context.Context, coin *entities.Coin) error {
	//TODO implement me
	panic("implement me")
}

func (c *CoinRepository) Delete(ctx context.Context, symbol string) error {
	_, err := c.db.ExecContext(ctx, `DELETE FROM coins WHERE symbol=$1`, symbol)
	return fmt.Errorf("failed delete coin: %w", err)
}

func NewCoinRepository(db *sql.DB) *CoinRepository {
	return &CoinRepository{db: db}
}
