package usecase

import (
	"golang-api-film-management/internal/domain"
)

// TODO: enhance: add logger to all the package (this is to check each files of this package)

type UserService interface {
	CreateUser(user domain.User) (*domain.User, error)
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
func (u *UserServiceImpl) CreateUser(user domain.User) (*domain.User, error) {
	err := u.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsernameAndPassword fetches a user from the user repository by their username and password.
func (u *UserServiceImpl) GetUserByUsernameAndPassword(userName, password string) (*domain.User, error) {
	foundUser, err := u.userRepo.GetUserByUsernameAndPassword(userName, password)
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// GetUserById retrieves a user by their primary key: user ID.
func (u *UserServiceImpl) GetUserById(userId int) (*domain.User, error) {
	foundUser, err := u.userRepo.GetUserById(userId)
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

// GetAllUsers returns a list of all users from the repository.
func (u *UserServiceImpl) GetAllUsers() (*[]domain.User, error) {
	users, err := u.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	return &users, nil
}
