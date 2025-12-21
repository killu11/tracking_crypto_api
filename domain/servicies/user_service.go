package servicies

import (
	"context"
	"crypto_api/domain/entities"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func (u UserRepository) Save(ctx context.Context, user *entities.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) FindByUsername(ctx context.Context, name string) (*entities.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) UpdatePassword(ctx context.Context, user *entities.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
