package domain

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByUsernameAndPassword(userName, password string) (User, error)
	GetUserById(userID int) (User, error)
	GetAllUsers() ([]User, error)
}
