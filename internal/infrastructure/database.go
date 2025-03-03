package infrastructure

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// FIXME: Should I think a better var name? gorm.DB should have a struct?
type Variable struct {
	Host             string
	Port             string
	UserName         string
	DatabaseName     string
	DatabasePassword string
	SSLMode          string
}

func OpenDatabase() (*gorm.DB, error) {
	variables := Variable{
		Host:             os.Getenv("DB_HOST"),
		Port:             os.Getenv("DB_PORT"),
		SSLMode:          os.Getenv("DB_SSLMODE"),
		UserName:         os.Getenv("POSTGRES_USER"),
		DatabaseName:     os.Getenv("POSTGRES_DB"),
		DatabasePassword: os.Getenv("POSTGRES_PASSWORD"),
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", variables.Host, variables.Port, variables.UserName, variables.DatabaseName, variables.DatabasePassword, variables.SSLMode)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Println("Fail to open the connection to database...")
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Fail to get the database instance...")
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		log.Println("Fail to verify the connection to database...")
		return nil, err
	}
	return db, nil
}

func CloseDatabase(db *gorm.DB) (*gorm.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Fail to get the database instance...")
		return nil, err
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Fail to close the database...")
		return nil, err
	}
	log.Println("Database connection is gracefully closed")
	return db, nil
}
