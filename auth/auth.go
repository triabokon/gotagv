package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var jwtKey = []byte("my_secret_key")

// CreateToken generates a new JWT token for a given username.
func CreateToken(username string) (string, error) {
	// Set token expiration time to be 30 minutes from now for testing purposes
	expirationTime := time.Now().Add(1 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a given token.
func ValidateToken(tknStr string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, err
	}

	return claims, nil
}

func HandleAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := ValidateToken(splitToken[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "username", claims.Username))

		next(w, r)
	}
}
