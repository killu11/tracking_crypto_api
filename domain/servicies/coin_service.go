package servicies

import (
	"context"
	"crypto_api/api_client/geckocoin"
	"crypto_api/domain/entities"
	repos "crypto_api/domain/repositories"
	"crypto_api/infrastructure/dto"
	cache "crypto_api/infrastructure/persistence/redis"
	"crypto_api/pkg"
	"errors"
	"fmt"
	"time"
)

type CoinService struct {
	coinRepo     repos.CoinRepository
	cache        *cache.CoinCache
	geckoCoin    *geckocoin.GeckoClient
	trackingRepo repos.TrackingRepository
}

var (
	ErrCoinNotFound        = errors.New("coin not found")
	ErrCoinAlreadyTracking = errors.New("coin already tracking")
	ErrNotTracking         = errors.New("coin doesn't tracking")
	ErrZeroTrackableCoins  = errors.New("zero trackable coins")
	ErrHasCeasedToExist    = errors.New("coin has ceased to exist")
)

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

// TrackCoin
// выполняет поиск монеты по символу;
// добавляет монету в отслеживание для соответствующего пользователя
func (s *CoinService) TrackCoin(
	ctx context.Context,
	symbol string,
	userID int,
) (*dto.TrackableCoinResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	coin, err := s.getOrCreateCoin(ctx, symbol)
	if err != nil {
		return nil, err
	}

	if err = s.trackingRepo.Add(ctx, userID, coin.Symbol); err != nil {
		if errors.Is(err, repos.ErrCoinAlreadyTracking) {
			return nil, ErrCoinAlreadyTracking
		}
		return nil, err
	}
	return dto.NewTrackableCoinResponse(coin), nil
}

// GetTrackableCoin
// Находит по символу монету из отслеживаемых пользователем
func (s *CoinService) GetTrackableCoin(
	ctx context.Context,
	userID int,
	symbol string,
) (*dto.TrackableCoinResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	coin, err := s.trackingRepo.FindBySymbol(ctx, userID, symbol)
	if err != nil {
		if errors.Is(err, repos.ErrCoinNotTracking) {
			return nil, ErrNotTracking
		}
		return nil, err
	}
	return dto.NewTrackableCoinResponse(coin), nil
}

// GetTrackableCoinsList
// Возвращает список отслеживаемых монет
func (s *CoinService) GetTrackableCoinsList(ctx context.Context, userID int) (*dto.TrackableListResponse, error) {
	list, err := s.trackingRepo.GetAll(ctx, userID)
	if err == nil {
		return &dto.TrackableListResponse{Cryptos: list}, nil
	}

	if errors.Is(err, repos.ErrZeroTrackableCoins) {
		return nil, ErrZeroTrackableCoins
	}
	return nil, err
}

func (s *CoinService) GetCoinStats(
	ctx context.Context,
	userID int,
	symbol string,
) (*dto.StatisticResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	stats, currentPrice, err := s.trackingRepo.GetStatsBySymbol(ctx, userID, symbol)

	if err != nil {
		if errors.Is(err, repos.ErrCoinNotTracking) {
			return nil, ErrNotTracking
		}
		return nil, err
	}

	stats.Change = stats.Max - currentPrice
	stats.PercentChange = stats.Change / (stats.Max / 100)
	return &dto.StatisticResponse{
		Symbol:         symbol,
		Price:          currentPrice,
		PriceStatistic: stats,
	}, nil
}

// GetPriceHistory
// Возвращает историю обновлений цены монеты
func (s *CoinService) GetPriceHistory(
	ctx context.Context,
	userID int,
	symbol string,
) (*dto.HistoryResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	history, err := s.trackingRepo.GetPriceHistory(ctx, userID, symbol)

	if err == nil {
		return &dto.HistoryResponse{
			Symbol:  symbol,
			History: history,
		}, nil
	}

	if errors.Is(err, repos.ErrCoinNotTracking) {
		return nil, ErrNotTracking
	}
	return nil, err
}

func (s *CoinService) RefreshTrackableCoin(
	ctx context.Context,
	userID int,
	symbol string,
) (*dto.TrackableCoinResponse, error) {
	coinID, found := s.cache.GetCryptoID(ctx, symbol)
	if found {
		coin, err := s.geckoCoin.GetCoinByID(ctx, coinID)
		if err != nil {
			return nil, err
		}
		return dto.NewTrackableCoinResponse(coin), nil
	}

	coinID, err := s.geckoCoin.GetCoinID(ctx, symbol)
	if err != nil {
		if errors.Is(err, repos.ErrCoinNotFound) {
			s.cache.SetNotFoundStatus(ctx, symbol)
			return nil, ErrCoinNotFound
		}
		return nil, err
	}

	coin, err := s.geckoCoin.GetCoinByID(ctx, coinID)
	if err != nil {
		if errors.Is(err, repos.ErrCoinNotFound) {
			s.cache.DropCryptoID(ctx, symbol)
			s.cache.SetNotFoundStatus(ctx, symbol)
			return nil, ErrCoinNotFound
		}
		return nil, err
	}

	if err = s.trackingRepo.UpdatePrice(ctx, coin, userID); err != nil {
		if errors.Is(err, repos.ErrCoinNotTracking) {
			return nil, fmt.Errorf("not rights to refresh coin price: %w", ErrNotTracking)
		}
		return nil, err

	}
	return dto.NewTrackableCoinResponse(coin), nil
}

// StopTrackCoin
// Останавливает отслеживание
func (s *CoinService) StopTrackCoin(
	ctx context.Context,
	userID int,
	symbol string,
) error {
	symbol = pkg.NormalizeSymbol(symbol)
	exists, err := s.trackingRepo.Exists(ctx, userID, symbol)
	if err != nil {
		return err
	}

	if !exists {
		return ErrNotTracking
	}
	return s.trackingRepo.Delete(ctx, userID, symbol)
}

// getOrCreateCoin
// Работает с настроенным API клиентом, кешом и БД
// Возвращает монету, найдя, либо сохранив ее в БД
func (s *CoinService) getOrCreateCoin(
	ctx context.Context,
	symbol string,
) (*entities.Coin, error) {
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
		s.cache.SetCryptoID(ctx, symbol, coinID, time.Hour*1)
	}

	coin, err = s.geckoCoin.GetCoinByID(ctx, coinID)
	if err != nil {
		if errors.Is(err, geckocoin.ErrCoinNotFound) {
			s.cache.SetNotFoundStatus(ctx, symbol)
			return nil, ErrCoinNotFound
		}
		return nil, err
	}

	if err = s.coinRepo.Save(ctx, coin); err != nil {
		return nil, err
	}
	return coin, nil
}

//func (s *CoinService) getCoinIDWithRetries(ctx context.Context, symbol string) (string, error) {
//	var lastErr error
//	const MaxAttempts = 3
//	for attempt := 1; attempt <= MaxAttempts; attempt++ {
//		coin, err := s.geckoCoin.GetCoinID(ctx, symbol)
//
//		if err == nil {
//			return coin, nil
//		}
//		lastErr = err
//		if errors.Is(err, ErrCoinNotFound) {
//			break
//		}
//
//		select {
//		case <-ctx.Done():
//			return "", ctx.Err()
//		case <-time.After(time.Duration(attempt*100) * time.Millisecond):
//			continue
//		}
//	}
//	return "", lastErr
//}
