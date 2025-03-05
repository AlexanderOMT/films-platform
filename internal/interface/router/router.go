package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang-api-film-management/internal/usecase"
)

type Server struct {
	*http.Server
}

func NewServerAtAddr(addr string, userService usecase.UserService, authService usecase.AuthService) *Server {
	RegisterRoutes(userService, authService)
	server := &http.Server{
		Addr: addr,
	}

	return &Server{server}
}

func (server *Server) Start() {
	// Start the server and listen for incoming connections
	log.Printf("Server is listening at address: %v", server.Addr)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server could not start: %v", err)
		}
	}()
	server.gracefullyShutdown()
}

func (server *Server) gracefullyShutdown() {
	// Create channel to listen for shutdown signal from OS or docker signal (note: see why is not working for 'ctl + cancel' in docker compose up)
	stopGracefulChannel := make(chan os.Signal, 1)
	signal.Notify(stopGracefulChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	log.Print("Server is waiting for shutdown signal from OS...")
	<-stopGracefulChannel
	log.Print("Server shutdown signal received")

	// Create a context with timeout to shutdown the opened connections (e.g: database connection)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Ensure to release resources
	defer cancel()

	log.Print("Shutting down the server")
	// Shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server failed to shutdown gracefully: %v", err)
	}
	log.Print("Server is gracefully shutdown")
}
