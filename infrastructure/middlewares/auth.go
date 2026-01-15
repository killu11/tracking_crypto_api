package middlewares

import (
	"crypto_api/pkg"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authToken := request.Header.Get("Authorization")
		if authToken == "" {
			pkg.JSONError(writer, "not found authorization token", http.StatusUnauthorized)
			return
		}
		parts := strings.Split(authToken, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			pkg.JSONError(writer, "uncorrected authorization header", http.StatusUnauthorized)
			return
		}
		// TODO: Закончить middleware,реализовав в pkg JWT модуль
		next.ServeHTTP(writer, request)
	})
}
