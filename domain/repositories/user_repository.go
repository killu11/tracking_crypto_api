package repositories

import (
	"context"
	"crypto_api/domain/entities"
)

type UserRepository interface {
	Save(ctx context.Context, user *entities.User) error
	FindByUsername(ctx context.Context, name string) (*entities.User, error)
	UpdatePassword(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id int) error
}
