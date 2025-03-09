package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

// AuthenticateTokenUser acts as a middleware that validates the jwt token in the request header.
// It checks if the token is valid and decodes the claims to extract the user subject ID.
// If can parse succesfully the token and get the subject ID, then will add the subject id to the request context and continious with the next handler
// If the token is invalid or malformed, it returns an error response.
func AuthenticateTokenUser(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing authorization token in the header", http.StatusBadRequest)
			return
		}

		tokenParsed, err := jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if !tokenParsed.Valid {
			log.Println("Error validating the jwt in the middleware: Token is invalid")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			log.Printf("Error validating the jwt in the middleware: %b", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
		subjectId := int(claims["Subject"].(float64))

		ctx := context.WithValue(r.Context(), "subjectId", subjectId)
		nextHandlerFunc.ServeHTTP(w, r.WithContext(ctx))
	}

}
