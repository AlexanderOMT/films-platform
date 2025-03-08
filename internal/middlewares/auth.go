package middlewares

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func AuthenticateTokenUser(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	// FIXME: style: status should be related to the actual error
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body and decode json
		tokenString := r.Header.Get("Authorization")
		tokenVerified, err := jwt.Parse(tokenString, func(tokenString *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if !tokenVerified.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		claims, ok := tokenVerified.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		}
		subjectId := int(claims["Subject"].(float64))

		// Create a context with the  user id for the handlers can make queries easier without parsing the token each time to DB
		ctx := context.WithValue(r.Context(), "subjectId", subjectId)
		nextHandlerFunc.ServeHTTP(w, r.WithContext(ctx))
	}

}
