package handlers

import (
	"crypto_api/application/services"
	"crypto_api/infrastructure/dto/request"
	"crypto_api/pkg"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CoinHandler struct {
	service *services.CoinService
}

func NewCoinHandler(s *services.CoinService) *CoinHandler {
	return &CoinHandler{service: s}
}

func (h *CoinHandler) InitHandlers(mux *chi.Mux) {
	mux.Post("/crypto", h.StartTracking)
}

// StartTracking
// Ручка для трекинга монеты
// В теле запроса symbol криптовалюты
func (h *CoinHandler) StartTracking(w http.ResponseWriter, r *http.Request) {
	// TODO: Добавить извлечение id пользователя из контекста
	var req request.TrackingCoinRequest
	if json.NewDecoder(r.Body).Decode(&req) != nil {
		pkg.JSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	// TODO: Заменить 1 на id пользователя
	coin, err := h.service.TrackCoin(r.Context(), req.Symbol, 1)
	if err != nil {
		h.matchError(w, err)
		return
	}

	pkg.JSONResponse(w, coin, http.StatusCreated)
}

func (h *CoinHandler) TrackableList(w http.ResponseWriter, r *http.Request) {
	// TODO: Добавить извлечение id пользователя из контекста
	response, err := h.service.GetTrackableCoinsList(r.Context(), 1)
	if err != nil {
		h.matchError(w, err)
		return
	}
	pkg.JSONResponse(w, response, http.StatusOK)
}

func (h *CoinHandler) matchError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(services.ErrCoinNotFound, err):
		pkg.JSONError(w, "coin not found", http.StatusNotFound)
	case errors.Is(services.ErrCoinAlreadyTracking, err):
		pkg.JSONError(w, err.Error(), http.StatusConflict)
	default:
		log.Printf("|WARNING|: %v", err)
		pkg.JSONError(w, "internal server error", http.StatusInternalServerError)
	}
}
