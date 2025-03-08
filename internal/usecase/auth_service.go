package usecase

import (
	"golang-api-film-management/internal/domain"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService interface {
	RegisterUser(user domain.User) (*domain.User, error)
	LoginUser(user domain.User) (*string, error)
	GetUserByToken(authenticatedToken *jwt.Token) (*domain.User, error)

	generateNewToken(user domain.User) (*string, error)
	authenticateToken(tokenToVerify string) (*jwt.Token, error)
}

type AuthServiceImpl struct {
	userService UserService
}

// FIXME: enhance in future: should be or auto-generated or given by env var
var secretKey = []byte("secret-key")

func NewAuthService(userService UserService) AuthService {
	return &AuthServiceImpl{userService: userService}
}

// Only add the user to DB, not generate any tokens
func (a *AuthServiceImpl) RegisterUser(user domain.User) (*domain.User, error) {
	createdUser, err := a.userService.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}

// Find out the user to DB by username and the passowrd. If there is a user for it, then generates a JWT token
func (a *AuthServiceImpl) LoginUser(user domain.User) (*string, error) {
	foundUser, err := a.userService.GetUserByUsernameAndPassword(user.Username, user.Password)
	if err != nil {
		log.Println("AuthServiceImpl: LoginUser user credentials are incorrect")
		return nil, err
	}
	log.Println("AuthServiceImpl: LoginUser user credentials are correct")

	createdToken, err := a.generateNewToken(*foundUser)
	if err != nil {
		log.Println("AuthServiceImpl: LoginUser error")
	}
	return createdToken, nil
}

// Find out the user to DB and authenticate the user token
// TODO: This is not used anymore, and it is implemented something of this in middleware because is more centralized. So remove this func when refactor user features
func (a *AuthServiceImpl) GetUserByToken(authenticatedToken *jwt.Token) (*domain.User, error) {
	claims := authenticatedToken.Claims.(jwt.MapClaims)

	subjectId := int(claims["Subject"].(float64))
	foundUser, err := a.userService.GetUserById(subjectId)
	if err != nil {
		log.Printf("Error extracting the user from the JWT: %v", subjectId)
		return nil, err
	}
	log.Printf("Token authenticated for subject Id: %v", subjectId)
	return foundUser, nil
}

func (a *AuthServiceImpl) generateNewToken(user domain.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"Subject":   user.Id,                          // Subject (user). Id is useful to query with the primary key of the database
			"ExpiresAt": time.Now().Add(time.Hour).Unix(), // expiration (time.Hour is a const of: 1h0m0s)
			"IssuesAt":  time.Now().Unix(),                // issued at
		})

	tokenSigned, err := token.SignedString(secretKey)
	if err != nil {
		return nil, err
	}
	return &tokenSigned, nil
}

func (a *AuthServiceImpl) authenticateToken(tokenToVerify string) (*jwt.Token, error) {
	tokenVerified, err := jwt.Parse(tokenToVerify, func(tokenString *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if !tokenVerified.Valid {
		log.Println("Token no valid")
		return nil, jwt.ErrSignatureInvalid
	}
	if err != nil {
		log.Printf("Token error %v", err)
		return nil, err
	}
	return tokenVerified, nil
}
