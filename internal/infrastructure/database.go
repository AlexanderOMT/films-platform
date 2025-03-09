package infrastructure

import (
	"context"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/sethvargo/go-envconfig"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// see https://www.reddit.com/r/golang/comments/uuet4u/loading_environment_variables_properly_in_go_with/ and https://github.com/sethvargo/go-envconfig
type DatabaseConfig struct {
	Host             string `env:"DB_HOST,required"`
	Port             string `env:"DB_PORT,required"`
	ConnSSLMode      string `env:"DB_SSLMODE,required"`
	UserName         string `env:"POSTGRES_USER,required"`
	DatabaseName     string `env:"POSTGRES_DB,required"`
	DatabasePassword string `env:"POSTGRES_PASSWORD,required"`
}

// OpenDatabase set a connection to the PostgreSQL database using credentials from environment variables.
// It returns a GORM database instance and any error if encountered.
func OpenDatabase() (*gorm.DB, error) {
	ctx := context.Background()
	var databaseConfigVar DatabaseConfig

	if err := envconfig.Process(ctx, &databaseConfigVar); err != nil {
		log.Printf("Failed opening the dababase connection: error loading and validating the configuration variables from enviroment")
		return nil, err
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", databaseConfigVar.Host, databaseConfigVar.Port, databaseConfigVar.UserName, databaseConfigVar.DatabaseName, databaseConfigVar.DatabasePassword, databaseConfigVar.ConnSSLMode)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Println("Fail to open the connection to database...")
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Fail to retrieve the database instance...")
		return nil, err
	}
	if err := sqlDB.Ping(); err != nil {
		log.Println("Fail to verify the connection to database...")
		return nil, err
	}
	return db, nil
}

// CloseDatabase closes the open connection to the PostgreSQL database and prevents new queries from starting. Close then waits for all queries that have started processing on the server to finish.
// It returns the database closed, or any error if encountered.
func CloseDatabase(db *gorm.DB) (*gorm.DB, error) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Fail to retrieve the database instance...")
		return nil, err
	}
	if err := sqlDB.Close(); err != nil {
		log.Println("Fail to close the database...")
		return nil, err
	}
	log.Println("Database connection is gracefully closed")
	return db, nil
}
