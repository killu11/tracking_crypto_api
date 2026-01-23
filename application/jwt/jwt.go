package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Tokens struct {
}

type TokenManager interface {
	NewJWT(userID string, ttl time.Duration) (string, error)
	Parse(token string) (string, error)
	NewRefreshToken() (string, error)
}

type Manager struct {
	secretKey string
}

func NewManager(secret string) (*Manager, error) {
	if secret == "" {
		return nil, fmt.Errorf("empty secret key")
	}
	return &Manager{secretKey: secret}, nil
}

func (m *Manager) NewJWT(userID string, ttl time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
	})

	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) Parse(accessToken string) (string, error) {
	var claims jwt.RegisteredClaims
	_, err := jwt.ParseWithClaims(accessToken, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexcepted signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretKey), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed parse jwt: %w", err)
	}
	if claims.Subject == "" {
		return "", fmt.Errorf("unexpected jwt payload")
	}
	return claims.Subject, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed generate refresh token: %v", err)
	}
	return hex.EncodeToString(buf), nil
}
