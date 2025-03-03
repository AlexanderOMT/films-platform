package usecase

import (
	"golang-api-film-management/internal/domain"
)

// TODO: enhance logger

type UserService interface {
	CreateUser(user domain.User) (*domain.User, error)
}

type UserServiceImpl struct {
	userRepo domain.UserRepository
}

func NewUserService(userRepo domain.UserRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

func (u *UserServiceImpl) CreateUser(user domain.User) (*domain.User, error) {
	err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
