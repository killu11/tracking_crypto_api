package entities

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidPassword = errors.New("invalid_password")
)

type User struct {
	Username string
	password []byte
}

func NewUser(username string) *User {
	return &User{Username: username}
}

func (u *User) SetPassword(password string) error {
	if len(password) <= 8 {
		return ErrInvalidPassword
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt generation password: %w", err)
	}

	u.password = pass
	return nil
}

func (u *User) PasswordHash() []byte {
	return u.password
}

// TODO: передвинуть функцию в транспортный слой
func UserIDFromContext(ctx context.Context) (int, bool) {
	id, ok := ctx.Value("userIDKey").(int)
	return id, ok
}
