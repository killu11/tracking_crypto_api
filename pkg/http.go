package pkg

import (
	"bytes"
	"crypto_api/infrastructure/dto/response"
	"encoding/json"
	"log"
	"net/http"
)

// JSONResponse - обрабатывает полученную структуру и кодирует её JSON-ом в responseWriter
func JSONResponse(w http.ResponseWriter, v any, code int) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(&v); err != nil {
		log.Printf("|WARNING| failed encode json: %v\n", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	buf.WriteTo(w)
}

// JSONError - обрабатывает полученное сообщение об ошибке и кодирует его JSON-ом в responseWriter
func JSONError(w http.ResponseWriter, msg string, code int) {
	errResponse := response.NewError(msg)
	JSONResponse(w, &errResponse, code)
}
