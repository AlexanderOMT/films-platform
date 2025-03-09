package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang-api-film-management/internal/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	RegisterUser(user *domain.User) (*domain.User, error)
	LoginUser(user domain.User) (*string, error)

	generateNewToken(user domain.User) (*string, error)
}

type AuthServiceImpl struct {
	userService UserService
}

// FIXME: enhance in future: should be or auto-generated or given by env var
var secretKey = []byte("secret-key")

// NewAuthService creates a new instance of AuthService with the provided UserService.
// It returns an implementation of the AuthService interface.
func NewAuthService(userService UserService) AuthService {
	return &AuthServiceImpl{userService: userService}
}

// RegisterUser creates a new user. Only add the user to DB, not generate any tokens.
func (a *AuthServiceImpl) RegisterUser(user *domain.User) (*domain.User, error) {
	hashedPassword, err := a.generateHash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a hash password during user registration : %w", err)
	}

	user.Password = hashedPassword
	createdUser, err := a.userService.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

// LoginUser authenticates a user by verifying their username and password.
// If the credentials are correct, it generates a new JWT token for the user.
func (a *AuthServiceImpl) LoginUser(user domain.User) (*string, error) {
	hashedPassword, err := a.generateHash(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate a hash password during user login validation : %w", err)
	}

	foundUser, err := a.userService.GetUserByUsernameAndPassword(user.Username, hashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to login user %s: %w", user.Username, err)
	}

	createdToken, err := a.generateNewToken(*foundUser)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token for user %s: %w", user.Username, err)
	}

	return createdToken, nil
}

// generateNewToken generates a new JWT token for the given user.
// The token contains the user ID as the subject, an expiration time, and the issued at time set to the current time.
// It returns the signed token as a string pointer or an error if the token signing fails.
func (a *AuthServiceImpl) generateNewToken(user domain.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Subject":   user.Id,                          // Subject (user). Id is useful to query with the primary key of the database
			"ExpiresAt": time.Now().Add(time.Hour).Unix(), // expiration (time.Hour is a const of: 1h0m0s)
			"IssuedAt":  time.Now().Unix(),                // issued at
		})

	tokenSigned, err := token.SignedString(secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign JWT token for user ID %d: %w", user.Id, err)
	}

	return &tokenSigned, nil
}

// generateHash generates the hash using `sha256` without salt for the given text.
// It returns the hashed text or error if any error encountred
func (a *AuthServiceImpl) generateHash(text string) (string, error) {
	algorithm := sha256.New()
	_, err := algorithm.Write([]byte(text))
	if err != nil {
		return "", fmt.Errorf("error generating hash for the given input text")
	}
	return hex.EncodeToString(algorithm.Sum(nil)), nil
}
