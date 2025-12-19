package postgres

import (
	"context"
	"crypto_api/domain/repositories"
	"database/sql"
	"errors"
)

type TrackingRepository struct {
	db *sql.DB
}

func NewTrackingRepository(db *sql.DB) *TrackingRepository {
	return &TrackingRepository{db: db}
}

func (t *TrackingRepository) Add(ctx context.Context, userID int, symbol string) error {
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
		return repositories.ErrAlreadyTracking
	}
	return nil
}

func (t *TrackingRepository) Exists(ctx context.Context, userID int, symbol string) (bool, error) {
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

func (t *TrackingRepository) Delete(ctx context.Context, userID int, symbol string) error {
	panic("Implement me")
}
