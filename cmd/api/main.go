package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"golang-api-film-management/internal/infrastructure"
	"golang-api-film-management/internal/interface/router"
	"golang-api-film-management/internal/usecase"
)

func main() {
	// Connect to database
	dbConnection, err := infrastructure.OpenDatabase()
	if err != nil {
		log.Fatalf("Fail to connect to database: %v", err)
	}
	log.Printf("Server connection with Database is succesfully established")

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	defer infrastructure.CloseDatabase(dbConnection)

	// Get the Database implementation
	userRepo := infrastructure.NewUserRepo(dbConnection)
	filmRepo := infrastructure.NewFilmRepo(dbConnection)
	// Get the services implementation for the DB implementation
	userService := usecase.NewUserService(userRepo)
	authService := usecase.NewAuthService(userService)
	filmService := usecase.NewFilmService(filmRepo)

	// Create a HTTP server with the services
	httpServer := router.NewServer(userService, authService, filmService)

	// Channel for stop signal
	httpServer.Start()

	<-quitChan
}
