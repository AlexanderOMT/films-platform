package router

import (
	"net/http"

	"golang-api-film-management/internal/interface/controller"
	"golang-api-film-management/internal/middlewares"
	"golang-api-film-management/internal/usecase"
)

func RegisterRoutes(userService usecase.UserService, authService usecase.AuthService) {
	// TODO: complete all the routes
	// TODO: enhance film service
	// TODO: enhance: sanitize the input (mitigate vulnerabilities)

	userHandler := controller.NewUserHandler(userService)
	authHandler := controller.NewAuthHandler(authService)

	http.HandleFunc("GET /users", userHandler.GetUsers)
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			authHandler.LoginUser(w, r)
		case http.MethodPost:
			authHandler.RegisterUser(w, r)
		case http.MethodDelete:
			middlewares.AuthenticateTokenUser(authHandler.DeleteUser)(w, r)
		default:
			http.Error(w, "Method for this route is not allowed", http.StatusMethodNotAllowed)
		}
	})

}
