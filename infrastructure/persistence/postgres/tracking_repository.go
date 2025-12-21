package postgres

import (
	"context"
	"crypto_api/domain/entities"
	"crypto_api/domain/repositories"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

type TrackingRepository struct {
	db *sql.DB
}

func NewTrackingRepository(db *sql.DB) *TrackingRepository {
	return &TrackingRepository{db: db}
}

func (t *TrackingRepository) Add(
	ctx context.Context,
	userID int,
	symbol string,
) error {
	res, err := t.db.ExecContext(
		ctx,
		`INSERT INTO users_tracking_coins(user_id, coin_symbol) 
				VALUES ($1,$2) ON CONFLICT DO NOTHING`,
		userID,
		symbol,
	)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()

	if err != nil {
		return err
	}
	if count == 0 {
		return repositories.ErrCoinAlreadyTracking
	}
	return nil
}

func (t *TrackingRepository) FindBySymbol(
	ctx context.Context,
	userID int,
	symbol string,
) (*entities.Coin, error) {
	var coin entities.Coin
	err := t.db.QueryRowContext(
		ctx,
		`SELECT symbol, name, current_price, last_updated FROM coins
	JOIN users_tracking_coins as utc ON utc.coin_symbol = coins.symbol
	WHERE utc.user_id = $1 AND utc.coin_symbol = $2`,
		userID,
		symbol,
	).Scan(coin.Symbol, coin.Name, coin.Usd, coin.LastUpdateAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrCoinNotTracking
		}
		return nil, err
	}
	return &coin, nil
}

func (t *TrackingRepository) GetStatsBySymbol(
	ctx context.Context,
	userID int,
	symbol string,
) (*entities.PriceStatistic, float64, error) {
	var stats entities.PriceStatistic
	var price float64
	err := t.db.QueryRowContext(
		ctx,
		`SELECT min(price) as min_price, max(price) as max,
       avg(price) as avg, count(*) as records,
       (
           SELECT price FROM coins_price_history
           WHERE symbol=$1
           ORDER BY timestamp DESC
           LIMIT 1
       ) as last_price
		FROM coins_price_history
		WHERE symbol = $2
	  	AND EXISTS(SELECT 1 FROM (
		   SELECT 1 FROM users_tracking_coins
		   WHERE coin_symbol=$3
			 AND user_id=$4
	   ) as utc
) `,
		symbol,
		symbol,
		symbol,
		userID,
	).Scan(&stats.Min, &stats.Max, &stats.Avg, &stats.Records, &price)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, 0, repositories.ErrCoinNotTracking
		}
		return nil, 0, err
	}

	return &stats, price, nil
}

func (t *TrackingRepository) GetPriceHistory(
	ctx context.Context,
	userID int,
	symbol string,
) ([]*entities.Price, error) {
	var resBytes []byte
	err := t.db.QueryRowContext(
		ctx,
		`SELECT COALESCE (json_agg(
				   json_build_object(
					   'price', price,
					   'timestamp', timestamp
				   )
    		  ),
       		'[]'::json
       ) FROM coins_price_history as h
    		   WHERE h.symbol=$2
    		   AND EXISTS(SELECT 1 FROM users_tracking_coins 
    		   WHERE user_id=$1 
    		   AND coin_symbol=$2) 
    		   GROUP BY h.timestamp`,
		userID,
		symbol).Scan(&resBytes)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repositories.ErrCoinNotTracking
		}
		return nil, err
	}

	var history []*entities.Price
	if err = json.Unmarshal(resBytes, &history); err != nil {
		return nil, fmt.Errorf("get history return's invalid json: %w", err)
	}

	//if len(history) == 0 {
	//	return nil, repositories.ErrCoinNotTracking
	//}
	return history, nil
}

func (t *TrackingRepository) Exists(
	ctx context.Context,
	userID int,
	symbol string,
) (bool, error) {
	var exists int
	err := t.db.QueryRowContext(
		ctx,
		`SELECT 1 FROM users_tracking_coins
				WHERE $1 = user_id AND $2 = coin_symbol`,
		userID,
		symbol,
	).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists == 1, nil
}

func (t *TrackingRepository) Delete(
	ctx context.Context,
	userID int,
	symbol string,
) error {
	_, err := t.db.ExecContext(
		ctx,
		`DELETE FROM users_tracking_coins 
       WHERE coin_symbol=$1 AND user_id=$2`,
		symbol,
		userID,
	)
	return err
}
