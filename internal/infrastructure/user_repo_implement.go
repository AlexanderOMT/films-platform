package infrastructure

// TODO change name of database package

import (
	"golang-api-film-management/internal/domain"
	"log"

	_ "github.com/lib/pq" // SQL driver
	"gorm.io/gorm"
)

type UserRepo struct {
	dbConnection *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{dbConnection: db}
}

// Register user
func (pg *UserRepo) CreateUser(user domain.User) error {
	if err := pg.dbConnection.Create(&user).Error; err != nil {
		log.Println("User creatation failed")
		return err
	}
	return nil
}

// Login user  # TODO: enhance in future: must return JWT Token -> (JWT, error)
func (pg *UserRepo) ValidateUser(name, password string) error {
	return nil
}

// Delete user
func (pg *UserRepo) DeleteUser(userID int) error {
	return nil
}
