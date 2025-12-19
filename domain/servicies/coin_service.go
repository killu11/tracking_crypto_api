package servicies

import (
	"context"
	"crypto_api/api_client/geckocoin"
	"crypto_api/domain/entities"
	repos "crypto_api/domain/repositories"
	cache "crypto_api/persistence/redis"
	"crypto_api/pkg"
	"errors"
	"fmt"
)

type CoinService struct {
	coinRepo     repos.CoinRepository
	cache        *cache.CoinCache
	geckoCoin    *geckocoin.GeckoClient
	trackingRepo repos.TrackingRepository
}

var (
	ErrCoinNotFound = errors.New("coin not found")
)

func (s *CoinService) TrackCoin(ctx context.Context, symbol string, userID int) error {
	coin, err := s.getOrCreateCoin(ctx, symbol)
	if err != nil {
		return err
	}
	exists, err := s.trackingRepo.Exists(ctx, userID, coin.ID)

	if err != nil {
		return err
	}

	if exists {
		return repos.ErrAlreadyTracking
	}
	if err = s.trackingRepo.Add(ctx, userID, coin.Symbol); err != nil {
		return err
	}
	return nil
}

func NewCoinService(
	repository repos.CoinRepository,
	cache *cache.CoinCache,
	client *geckocoin.GeckoClient,
	trackingRepository repos.TrackingRepository,
) *CoinService {
	return &CoinService{
		coinRepo:     repository,
		cache:        cache,
		geckoCoin:    client,
		trackingRepo: trackingRepository,
	}
}

func (s *CoinService) getOrCreateCoin(ctx context.Context, symbol string) (*entities.Coin, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	coin, err := s.coinRepo.FindBySymbol(ctx, symbol)

	if err == nil {
		return coin, nil
	}

	if !errors.Is(err, repos.ErrCoinNotFound) {
		return nil, err
	}

	if s.cache.IsNotFound(ctx, symbol) {
		return nil, ErrCoinNotFound
	}

	coinID, found := s.cache.GetCryptoID(ctx, symbol)

	if !found {
		coinID, err = s.geckoCoin.GetCoinID(ctx, symbol)
		if err != nil {
			if errors.Is(err, geckocoin.ErrCoinNotFound) {
				s.cache.SetNotFoundStatus(ctx, symbol)
				return nil, ErrCoinNotFound
			}
			return nil, err
		}
	}

	coin, err = s.geckoCoin.GetCoinByID(ctx, coinID)
	if err != nil {
		if errors.Is(err, geckocoin.ErrCoinNotFound) {
			s.cache.SetNotFoundStatus(ctx, symbol)
			return nil, fmt.Errorf("coin not found")
		}
		return nil, err
	}

	if err = s.coinRepo.Save(ctx, coin); err != nil {
		return nil, err
	}
	return coin, nil
}
