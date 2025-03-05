package middlewares

import (
	"log"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret-key")

func AuthenticateTokenUser(nextHandlerFunc http.HandlerFunc) http.HandlerFunc {
	log.Printf("middleware:")
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse request body and decode json
		tokenToVerify := r.Header.Get("Authorization")
		tokenVerified, err := jwt.Parse(tokenToVerify, func(tokenString *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})
		if !tokenVerified.Valid {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest) // FIXME: style: status should be related to incorrect credentials ? bad request format?
			return
		}
		nextHandlerFunc(w, r)
	}

}
