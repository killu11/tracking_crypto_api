package servicies

import (
	"context"
	"crypto_api/api_client/geckocoin"
	"crypto_api/domain/entities"
	repos "crypto_api/domain/repositories"
	cache "crypto_api/persistence/redis"
	"errors"
	"fmt"
)

type CoinService struct {
	coinRepo     repos.CoinRepository
	cache        *cache.CoinCache
	geckoCoin    *geckocoin.GeckoClient
	trackingRepo repos.TrackingRepository
}

func (s *CoinService) TrackCoin(ctx context.Context, symbol string, userID int) error {
	coin, err := s.getOrCreateCoin(ctx, symbol)
	if err != nil {
		return err
	}
	exists, err := s.trackingRepo.Exists(ctx, userID, symbol)

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

func NewCoinService(repository repos.CoinRepository) *CoinService {
	return &CoinService{coinRepo: repository}
}

func (s *CoinService) getOrCreateCoin(ctx context.Context, symbol string) (*entities.Coin, error) {
	coin, err := s.coinRepo.FindBySymbol(ctx, symbol)

	if err == nil {
		return coin, nil
	}

	if !errors.Is(err, repos.ErrCoinNotFound) {
		return nil, err
	}

	ok, _ := s.cache.IsNotFound(ctx, symbol)
	if ok {
		//TODO: нужна сервисная ошибка
		return nil, fmt.Errorf("coin not found")
	}

	coinID, found, _ := s.cache.GetCryptoID(ctx, symbol)

	if !found {
		coinID, err = s.geckoCoin.GetCoinID(ctx, symbol)
		if err != nil {
			if errors.Is(err, geckocoin.ErrCoinNotFound) {
				s.cache.SetNotFoundStatus(ctx, symbol)
				//TODO: нужна сервисная ошибка
				return nil, fmt.Errorf("coin not found")
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
