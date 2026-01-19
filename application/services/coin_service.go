package services

import (
	"context"
	"crypto_api/api_client/geckocoin"
	"crypto_api/domain/entities"
	repos "crypto_api/domain/repositories"
	"crypto_api/infrastructure/dto/response"
	cache "crypto_api/infrastructure/persistence/redis"
	"crypto_api/pkg"
	"errors"
	"log"
	"time"
)

type CoinService struct {
	coinRepo     repos.CoinRepository
	trackingRepo repos.TrackingRepository
	cache        *cache.CoinCache
	geckoCoin    *geckocoin.GeckoClient
}

var (
	ErrCoinNotFound        = errors.New("coin not found")
	ErrCoinAlreadyTracking = errors.New("coin already tracking")
	ErrCoinNotTracking     = errors.New("coin doesn't tracking")
	ErrZeroTrackableCoins  = errors.New("zero trackable coins")
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
) (*response.TrackableCoinResponse, error) {
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
	return response.NewTrackableCoinResponse(coin), nil
}

// GetTrackableCoin
// Находит по символу монету из отслеживаемых пользователем
func (s *CoinService) GetTrackableCoin(
	ctx context.Context,
	userID int,
	symbol string,
) (*response.TrackableCoinResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	coin, err := s.trackingRepo.FindBySymbol(ctx, userID, symbol)
	if err != nil {
		if errors.Is(err, repos.ErrCoinNotTracking) {
			return nil, ErrCoinNotTracking
		}
		return nil, err
	}
	return response.NewTrackableCoinResponse(coin), nil
}

// GetTrackableCoinsList
// Возвращает список отслеживаемых монет
func (s *CoinService) GetTrackableCoinsList(
	ctx context.Context,
	userID int,
) (*response.TrackableListResponse, error) {
	list, err := s.trackingRepo.GetAll(ctx, userID)
	if err == nil {
		return &response.TrackableListResponse{Cryptos: list}, nil
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
) (*response.StatisticResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	stats, currentPrice, err := s.trackingRepo.GetStatsBySymbol(ctx, userID, symbol)

	if err != nil {
		if errors.Is(err, repos.ErrCoinNotTracking) {
			return nil, ErrCoinNotTracking
		}
		return nil, err
	}

	stats.Change = stats.Max - currentPrice
	stats.PercentChange = stats.Change / (stats.Max / 100)
	return &response.StatisticResponse{
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
) (*response.HistoryResponse, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	history, err := s.trackingRepo.GetPriceHistory(ctx, userID, symbol)

	if err == nil {
		return &response.HistoryResponse{
			Symbol:  symbol,
			History: history,
		}, nil
	}

	if errors.Is(err, repos.ErrCoinNotTracking) {
		return nil, ErrCoinNotTracking
	}
	return nil, err
}

func (s *CoinService) RefreshTrackableCoin(
	ctx context.Context,
	userID int,
	symbol string,
) (
	*response.TrackableCoinResponse,
	error,
) {
	symbol = pkg.NormalizeSymbol(symbol)

	var id string
	var err error
	const maxAttempts = 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		id, err = s.getIdBySymbol(ctx, symbol)
		if err == nil {
			break
		}

		if errors.Is(err, ErrCoinNotFound) {
			return nil, ErrCoinNotFound
		}

		if attempt == maxAttempts {
			break
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(100 * time.Duration(attempt) * time.Millisecond):
			continue
		}
	}

	freshCoinData, err := s.geckoCoin.GetFreshCoinData(ctx, id)
	if err != nil {
		if errors.Is(err, geckocoin.ErrCoinNotFound) {
			return nil, ErrCoinNotFound
		}
		return nil, err
	}

	coin := freshCoinData.ToEntity(symbol)
	log.Printf("Refresh this coin -> %v\n", coin)
	if err = s.trackingRepo.UpdatePrice(ctx, coin, userID); err != nil {
		if errors.Is(err, repos.ErrCoinNotTracking) {
			return nil, ErrCoinNotTracking
		}
		return nil, err
	}
	return response.NewTrackableCoinResponse(coin), nil
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
		return ErrCoinNotTracking
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

	coinID, err := s.getIdBySymbol(ctx, symbol)
	if err != nil {
		return nil, err
	}

	coin, err = s.geckoCoin.GetCoinByID(ctx, coinID)
	if err != nil {
		if errors.Is(err, geckocoin.ErrCoinNotFound) {
			s.cache.SetNotFoundStatus(ctx, symbol, 1*time.Hour)
			return nil, ErrCoinNotFound
		}
		return nil, err
	}

	if err = s.coinRepo.Save(ctx, coin); err != nil {
		return nil, err
	}
	return coin, nil
}

// getIdBySymbol поиск coinID среди кеша и/или внешнего
func (s *CoinService) getIdBySymbol(ctx context.Context, symbol string) (string, error) {
	symbol = pkg.NormalizeSymbol(symbol)
	if s.cache.IsNotFound(ctx, symbol) {
		return "", ErrCoinNotFound
	}

	if id, found := s.cache.GetCryptoID(ctx, symbol); found {
		return id, nil
	}

	id, err := s.geckoCoin.SymbolToID(ctx, symbol)
	if err == nil {
		s.cache.SetCryptoID(ctx, symbol, id, 1*time.Hour)
		return id, nil
	}

	if errors.Is(geckocoin.ErrCoinNotFound, err) {
		s.cache.SetNotFoundStatus(ctx, symbol, 1*time.Hour)
		return "", ErrCoinNotFound
	}

	return "", err
}
