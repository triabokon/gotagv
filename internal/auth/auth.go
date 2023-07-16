package auth

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKeys string

const UserIDKey ContextKeys = "user_id"

type Auth struct {
	config *Config
}

func New(config *Config) *Auth {
	srv := &Auth{
		config: config,
	}
	return srv
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// CreateToken generates a new JWT token for a given username.
func (a *Auth) CreateToken(userID string) (string, error) {
	// Set token expiration time to be 30 minutes from now for testing purposes
	expirationTime := time.Now().Add(1 * time.Minute)

	// Create the JWT claims, which includes the userID and expiry time
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(a.config.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken parses and validates a given token.
func (a *Auth) ValidateToken(tknStr string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.config.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, err
	}

	return claims, nil
}

func (a *Auth) HandleAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := a.ValidateToken(splitToken[1])
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), UserIDKey, claims.UserID))

		next(w, r)
	}
}