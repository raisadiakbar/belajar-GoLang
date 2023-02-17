package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for the given username
func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type response struct {
	Message string `json:"message"`
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	resp := response{Message: message}
	jsonResp, _ := json.Marshal(resp)

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResp)
}

// Authenticate middleware validates the JWT token in the request header
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			RespondWithError(w, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})

		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				RespondWithError(w, http.StatusUnauthorized, "Invalid token signature")
			} else {
				RespondWithError(w, http.StatusUnauthorized, "Invalid token")
			}
			return
		}

		if !token.Valid {
			RespondWithError(w, http.StatusUnauthorized, "Token is not valid")
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
