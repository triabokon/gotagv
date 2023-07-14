package server

import (
	"encoding/json"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/triabokon/gotagv/auth"
)

var users = map[string]bool{
	"user1": true,
	"user2": true,
}

type SignUpResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	userID := uuid.New()
	// todo: write user to db

	// Create a new token for the current use, and return it in the response
	tokenString, err := auth.CreateToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.JSONResponse(w, &SignUpResponse{UserID: userID, Token: tokenString})
}

type SignInRequest struct {
	UserID string `json:"user_id,omitempty"`
}

type SignInResponse struct {
	Token string `json:"token"`
}

func (s *Server) SignIn(w http.ResponseWriter, r *http.Request) {
	// Parse and decode the request body into a new `Credentials` instance
	req := &SignInRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		// If there is something wrong with the request body, return a 400 status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, ok := users[req.UserID] // todo: get user from db
	if !ok {
		// todo: change error from string to error
		s.ErrorResponse(w, "failed to check user", http.StatusUnauthorized)
		return
	}

	// Create a new token for the current use, and return it in the response
	tokenString, err := auth.CreateToken(req.UserID)
	if err != nil {
		s.ErrorResponse(w, "failed to create jwt token", http.StatusInternalServerError)
		return
	}

	s.JSONResponse(w, &SignInResponse{Token: tokenString})
}
