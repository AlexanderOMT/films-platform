package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"golang-api-film-management/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserService
}

func NewUserHandler(userService usecase.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers retrieves a list of all users.
// It calls the user service to fetch all users and returns the user list or an error.
func (userHandler *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userHandler.userService.GetAllUsers()
	if err != nil {
		log.Printf("Error retrieving users list: %v", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	log.Printf("Successfully retrieved all the user list | Users: %v", users)

	// Set response header to json type
	w.Header().Set("Content-Type", "application/json")
	// Encode users to json and write to response
	json.NewEncoder(w).Encode(users)
}
