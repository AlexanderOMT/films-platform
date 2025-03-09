package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang-api-film-management/internal/domain"
	"golang-api-film-management/internal/usecase"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService usecase.AuthService
}

func NewAuthHandler(authService usecase.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterUser registers a new user in the system.
// It decodes the user's name and password from the request body, and then calls the service to register the user.
// Its response is the user fields created or any error if encountred
func (authHandler *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var userToCreate domain.User
	if err := json.NewDecoder(r.Body).Decode(&userToCreate); err != nil {
		log.Printf("User creation error during mapping the body request to user model: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(userToCreate); err != nil {
		log.Printf("Missing required fields for user registration: %v", err)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	_, err := authHandler.authService.ValidateUserFields(&userToCreate)
	if err != nil {
		http.Error(w, fmt.Sprintf("User fields does not meet the required constraints: %v", err.Error()), http.StatusBadRequest)
		return
	}

	createdUser, err := authHandler.authService.RegisterUser(&userToCreate)
	if err != nil {
		log.Printf("User creation failed during registration: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	WriteJSONResponse(w, http.StatusOK, createdUser)
	log.Printf("User created successfully | Username: %v", createdUser.Username)
}

// LoginUser generate a new jwt token for the user.
// It call a service to validate the username and password given in the body.
// If this service validate the credentials successfully, then generates a new jwt token
// Its response is the new jwt token created or any error if encountred
func (a *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var userToLogin domain.User
	if err := json.NewDecoder(r.Body).Decode(&userToLogin); err != nil {
		log.Printf("User logging error during mapping the body request to user model: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(userToLogin); err != nil {
		log.Printf("Missing required fields for user login: %v", err)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	tokenStringUser, err := a.authService.LoginUser(userToLogin)
	if err != nil {
		log.Printf("User logging error during credentials verification: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	WriteJSONResponse(w, http.StatusOK, map[string]interface{}{"token": tokenStringUser})
	log.Printf("User logged successfully | Username: %v", userToLogin.Username)
}
