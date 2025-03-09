package usecase

import (
	"fmt"
	"golang-api-film-management/internal/domain"
)

type UserService interface {
	CreateUser(user *domain.User) (*domain.User, error)
	GetUserByUsernameAndPassword(userName, password string) (*domain.User, error)
	GetUserById(userId int) (*domain.User, error)
	GetAllUsers() (*[]domain.User, error)
}

type UserServiceImpl struct {
	userRepo domain.UserRepository
}

// NewUserService creates a new instance of UserService with the provided user repository.
func NewUserService(userRepo domain.UserRepository) UserService {
	return &UserServiceImpl{userRepo: userRepo}
}

// CreateUser saves a new user to the repository and returns the pointer of the created user.
func (u *UserServiceImpl) CreateUser(user *domain.User) (*domain.User, error) {
	err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

// GetUserByUsernameAndPassword fetches a user from the user repository by their username and password.
func (u *UserServiceImpl) GetUserByUsernameAndPassword(userName, password string) (*domain.User, error) {
	foundUser, err := u.userRepo.GetUserByUsernameAndPassword(userName, password)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username and password: %w", err)
	}
	return &foundUser, nil
}

// GetUserById retrieves a user by their primary key: user ID.
func (u *UserServiceImpl) GetUserById(userId int) (*domain.User, error) {
	foundUser, err := u.userRepo.GetUserById(userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", userId, err)
	}
	return &foundUser, nil
}

// GetAllUsers returns a list of all users from the repository.
func (u *UserServiceImpl) GetAllUsers() (*[]domain.User, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return &users, nil
}
