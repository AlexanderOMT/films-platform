package controller

import (
	"encoding/json"
	"net/http"

	"golang-api-film-management/internal/usecase"
)

type UserHandler struct {
	userService usecase.UserService
}

func NewUserHandler(userService usecase.UserService) *UserHandler {
	return &UserHandler{userService: userService}
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
