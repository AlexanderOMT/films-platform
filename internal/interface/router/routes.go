package router

import (
	"net/http"

	"golang-api-film-management/internal/interface/controller"
	"golang-api-film-management/internal/middlewares"
	"golang-api-film-management/internal/usecase"
)

func RegisterRoutes(userService usecase.UserService, authService usecase.AuthService, filmService usecase.FilmService) {

	// TODO: enhance: sanitize the input (to mitigate vulnerabilities e.g: sql injection)

	userHandler := controller.NewUserHandler(userService)
	authHandler := controller.NewAuthHandler(authService)
	filmHandler := controller.NewFilmHandler(filmService)

	// Related to user endpoints, some of them are unprotected routes (this is that does not require jwt)
	http.HandleFunc("GET /users", userHandler.GetUsers)
	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			authHandler.RegisterUser(w, r) // Missing username constraint and password validations
		case http.MethodGet:
			authHandler.LoginUser(w, r) // Done
		case http.MethodDelete:
			middlewares.AuthenticateTokenUser(authHandler.DeleteUser)(w, r) // No implemented, but no mandatory up to now
		default:
			http.Error(w, "Method for this route is not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Related to films endpoints, each route is a protected one
	http.HandleFunc("POST /film", middlewares.AuthenticateTokenUser(filmHandler.CreateFilm))   // Done
	http.HandleFunc("GET /films", middlewares.AuthenticateTokenUser(filmHandler.GetAllFilms))  // Done
	http.HandleFunc("PATCH /film", middlewares.AuthenticateTokenUser(filmHandler.PatchFilm))   // Done
	http.HandleFunc("PUT /film", middlewares.AuthenticateTokenUser(filmHandler.PutFilm))       // Done
	http.HandleFunc("DELETE /film", middlewares.AuthenticateTokenUser(filmHandler.DeleteFilm)) // Done
}
