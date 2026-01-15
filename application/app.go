package application

import (
	"crypto_api/api_client/geckocoin"
	"crypto_api/application/services"
	config2 "crypto_api/config"
	"crypto_api/domain/repositories"
	"crypto_api/infrastructure/handlers"
	"crypto_api/infrastructure/persistence/postgres"
	"crypto_api/infrastructure/persistence/redis"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type App struct {
	repos    Repositories
	services Services
	handlers Handlers
}

type Repositories struct {
	coinRepo     repositories.CoinRepository
	trackingRepo repositories.TrackingRepository
	coinCache    *redis.CoinCache
}
type Services struct {
	coinService *services.CoinService
}

type Handlers struct {
	coinHandler *handlers.CoinHandler
}

func (a *App) Start() {
	mux := chi.NewRouter()
	// TODO: Реализовать middleware для проверки JWT токена
	a.handlers.coinHandler.InitHandlers(mux)
	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadTimeout:       10 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
	log.Println("Server start on port :8080!")
	log.Fatalln(server.ListenAndServe())

}

func NewApp() *App {
	config, err := config2.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}
	// Инициализация апи-клиента
	geckoClient := geckocoin.NewGeckoClient(config.Gecko)
	// Ретрай пинга к стороннему апи
	for i := 1; i <= 3; i++ {
		err = geckoClient.Ping()
		if err == nil {
			break
		}
		if i == 3 {
			break
		}
		select {
		case <-time.After(100 * time.Duration(i) * time.Millisecond):
			continue
		}
	}

	if err != nil {
		log.Fatalln(err)
	}
	// Подключение к БД (Postgres, Redis)
	db, err := postgres.NewPostgresConnection(config.Postgres)
	if err != nil {
		log.Fatalf("failed app initing postgres connection: %v\n", err)
	}

	redisClient, err := redis.NewRedisConnection(config.Redis)
	if err != nil {
		log.Fatalln(err)
	}

	// Инициализация репозиториев
	repos := Repositories{
		coinRepo:     postgres.NewCoinRepository(db),
		trackingRepo: postgres.NewTrackingRepository(db),
		coinCache:    redis.NewCoinCacheRepository(redisClient),
	}

	// Инициализация сервисов
	s := Services{
		coinService: services.NewCoinService(
			repos.coinRepo,
			repos.coinCache,
			geckoClient,
			repos.trackingRepo,
		),
	}
	// Инициализация обработчиков
	h := Handlers{coinHandler: handlers.NewCoinHandler(s.coinService)}

	return &App{
		repos:    repos,
		services: s,
		handlers: h,
	}
}
