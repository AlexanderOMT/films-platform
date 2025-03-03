package domain

type UserRepository interface {
	CreateUser(user User) error
	DeleteUser(userID int) error
	ValidateUser(name, password string) error
}
