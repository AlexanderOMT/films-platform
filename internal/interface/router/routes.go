package router

import (
	"net/http"

	"golang-api-film-management/internal/interface/controller"
	"golang-api-film-management/internal/usecase"
)

func RegisterRoutes(userService usecase.UserService) {
	// TODO: complete all the routes
	// TODO: enhance film service

	userHandler := controller.NewUserHandler(userService)

	http.HandleFunc("/user", methodHandler("POST", userHandler.CreateUser))
}

func methodHandler(expectedMethod string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != expectedMethod {
			http.Error(w, "Method for this route is not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}
