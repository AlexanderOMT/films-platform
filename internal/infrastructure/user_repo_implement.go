package infrastructure

// TODO change name of database package

import (
	_ "github.com/lib/pq" // SQL driver
	"gorm.io/gorm"

	"golang-api-film-management/internal/domain"
)

type UserRepo struct {
	dbConnection *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{dbConnection: db}
}

// CreateUser creates a new given ser in the database connection
func (pg *UserRepo) CreateUser(user *domain.User) error {
	if err := pg.dbConnection.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// TODO: consider refactor: GetUserByUsername and GetUserById could be in one method: GetUserByKey(...)
func (pg *UserRepo) GetUserByUsernameAndPassword(userName, password string) (domain.User, error) {
	var user domain.User
	if err := pg.dbConnection.Where("username = ? AND password = ?", userName, password).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.User{}, err
		}
		return domain.User{}, err
	}
	return user, nil
}

// GetUserById gets a user for the given id
func (pg *UserRepo) GetUserById(userID int) (domain.User, error) {
	var user domain.User
	if err := pg.dbConnection.First(&user, userID).Error; err != nil {
		return domain.User{}, err
	}
	return user, nil
}

// GetAllUsers retrieves a list of all user records in the database
func (pg *UserRepo) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	if err := pg.dbConnection.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
