package router

import (
	"net/http"

	"golang-api-film-management/internal/interface/controller"
	"golang-api-film-management/internal/middlewares"
	"golang-api-film-management/internal/usecase"
)

func RegisterRoutes(userService usecase.UserService, authService usecase.AuthService, filmService usecase.FilmService) {

	userHandler := controller.NewUserHandler(userService)
	authHandler := controller.NewAuthHandler(authService)
	filmHandler := controller.NewFilmHandler(filmService)

	// Related to user endpoints, some of them are unprotected routes (this is that does not require jwt)
	http.HandleFunc("GET /users", userHandler.GetUsers)
	http.HandleFunc("POST /register", authHandler.RegisterUser) // Missing username constraint and password validations
	http.HandleFunc("POST /login", authHandler.LoginUser)

	// Related to films endpoints, each route is a protected one
	http.HandleFunc("POST /film", middlewares.AuthenticateTokenUser(filmHandler.CreateFilm))
	http.HandleFunc("GET /films", middlewares.AuthenticateTokenUser(filmHandler.GetAllFilms))
	http.HandleFunc("PATCH /film", middlewares.AuthenticateTokenUser(filmHandler.PatchFilm))
	http.HandleFunc("PUT /film", middlewares.AuthenticateTokenUser(filmHandler.PutFilm))
	http.HandleFunc("DELETE /film", middlewares.AuthenticateTokenUser(filmHandler.DeleteFilm))
}
