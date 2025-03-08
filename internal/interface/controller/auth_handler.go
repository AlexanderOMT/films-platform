package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"golang-api-film-management/internal/domain"
	"golang-api-film-management/internal/usecase"
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
	// TODO: enhance: 1) regex for alphanumeric starting with letter 2) validations (e.g: max length, and so on)

	// Parse request body and decode json
	var userToCreate domain.User
	if err := json.NewDecoder(r.Body).Decode(&userToCreate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: enhance: Hash for the user password before storing
	createdUser, err := authHandler.authService.RegisterUser(userToCreate)
	if err != nil {
		log.Println("handler: createUser error")
	}

	// FIXME: style: This block of code could be repeated from time to time, consider refactoring or creating a util function
	// Set response header to json type
	w.Header().Set("Content-Type", "application/json")
	// Encode users to json and write to response
	json.NewEncoder(w).Encode(userToCreate)

	log.Printf("Created User: %v", createdUser.Username)
}

// LoginUser generate a new jwt token for the user.
// It call a service to validate the username and password given in the body.
// If this service validate the credentials successfully, then generates a new jwt token
// Its response is the new jwt token created or any error if encountred
func (authHandler *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// TODO: rusername should be on URL??
	// Parse request body and decode json
	var userToLogin domain.User
	if err := json.NewDecoder(r.Body).Decode(&userToLogin); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tokenStringUser, err := authHandler.authService.LoginUser(userToLogin)
	if err != nil {
		// status should be related to incorrect credentials
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set response header to json type
	w.Header().Set("Content-Type", "application/json")
	// Encode users to json and write to response
	json.NewEncoder(w).Encode(tokenStringUser)

	log.Printf("Logged User token: %v", tokenStringUser)
}

func (authService *AuthHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// TODO: enhance: should only be executed by the itself. This is, only the user with his own token can deletes itself and no others users
	// TODO: username on the url?
	log.Println("Delete user endpoint")

	// Set response header to json type
	w.Header().Set("Content-Type", "application/json")
	// Encode users to json and write to response
	json.NewEncoder(w).Encode("Delete user endpoint")

}
