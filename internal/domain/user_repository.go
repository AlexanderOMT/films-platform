package domain

type UserRepository interface {
	CreateUser(user User) error
	GetUserById(userID int) (User, error)
	GetUserByUsernameAndPassword(userName, password string) (User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(userID int) error
}
