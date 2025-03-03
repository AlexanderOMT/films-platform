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
	// TODO: consider: Should we keep the connection up all the time? Or just we want to query to DB ?
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
	// Get the service implementation for the DB implementation
	userService := usecase.NewUserService(userRepo)

	// Create a HTTP server
	// TODO: addressToListen should be read from env or docker port
	addressToListen := "0.0.0.0:8000"
	httpServer := router.NewServerAtAddr(addressToListen, userService)

	// Channel for stop signal
	httpServer.Start()

	<-quitChan
	// os.Exit(0) This avoid the Database to be closed. If we open and close the database when is only querying to Postgres, then we could keep this
}
