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

func (userHandler *UserHandler) GetUserById(userId int) (*domain.User, error) {
	foundUser, err := userHandler.userService.GetUserById(userId)
	if err != nil {
		log.Printf("Error extracting the user from the JWT: %v", userId)
		return nil, err
	}
	return foundUser, nil
}

func (userHandler *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userHandler.userService.GetAllUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	// Set response header to json type
	w.Header().Set("Content-Type", "application/json")
	// Encode users to json and write to response
	json.NewEncoder(w).Encode(users)
}
