package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"golang-api-film-management/internal/domain"
	"golang-api-film-management/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserService
}

func NewUserHandler(userService usecase.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (userHandler *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	// Parse request body and decode json
	var userToCreate domain.User
	if err := json.NewDecoder(r.Body).Decode(&userToCreate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: enhance Hash for the user password before storing
	createdUser, err := userHandler.userService.CreateUser(userToCreate)
	if err != nil {
		log.Println("handler: createUser error")
	}

	// FIXME This block of code could be repeated from time to time, consider refactoring or creating a util function
	// Set response header to json type
	w.Header().Set("Content-Type", "application/json")
	// Encode users to json and write to response
	json.NewEncoder(w).Encode(userToCreate)

	log.Printf("Created Username: %v", createdUser.Username)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {

}

func getUser(w http.ResponseWriter, r *http.Request) {

}

func getUsers(w http.ResponseWriter, r *http.Request) {

}
