package router

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sethvargo/go-envconfig"

	"golang-api-film-management/internal/usecase"
)

type Server struct {
	*http.Server
}

type ServerConfig struct {
	Port string `env:"EXPOSE_PORT, default=8000"`
	Host string `env:"HOST, default=0.0.0.0"`
}

func NewServer(userService usecase.UserService, authService usecase.AuthService, filmService usecase.FilmService) *Server {
	ctx := context.Background()
	var databaseConfigVar ServerConfig

	if err := envconfig.Process(ctx, &databaseConfigVar); err != nil {
		log.Printf("Failed creating the server: error loading the configuration variables from enviroment")
	}

	RegisterRoutes(userService, authService, filmService)
	server := &http.Server{
		Addr: databaseConfigVar.Host + ":" + databaseConfigVar.Port,
	}
	return &Server{server}
}

// Start starts the server creating a go routine which will listen and server the connections, and set a gracefully shutdown which listen for shutdown signal from the system
func (server *Server) Start() {
	log.Printf("Listening at : %v", server.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server could not start: %v", err)
		}
	}()

	server.gracefullyShutdown()
}

// gracefullyShutdown creates a channel to listen for shutdown signal (Interrupt, SIGTERM, SIGINT, SIGQUIT). When receied the signal, this func will create a context with timeout to shutdown the opened connections (e.g: database connection) and shutdown the server gracefully
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
