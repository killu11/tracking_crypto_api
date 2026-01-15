package services

import repos "crypto_api/domain/repositories"

type UserService struct {
	repo repos.UserRepository
}
