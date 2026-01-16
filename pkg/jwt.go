package pkg

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenLife = 12
	issuer          = "auth-service"
)

type AccessTokenClaims struct {
	ID int `json:"user_id"`
	jwt.RegisteredClaims
}

func newAccessTokenClaims(id int) *AccessTokenClaims {
	return &AccessTokenClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   fmt.Sprintf("user-%d", id),
			Audience:  nil,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * accessTokenLife)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func CreateAccessTokenString(userID int) (string, error) {
	panic("implement me")
}

func NewRefreshToken() {
	panic("implement me")
}

func ValidatingAccessToken() {
	panic("implement me")
}

func getSecretKey() ([]byte, error) {
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		return nil, fmt.Errorf("not found secret key from env")
	}

	return []byte(secret), nil
}
