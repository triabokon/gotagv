package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pborman/uuid"

	"github.com/triabokon/gotagv/internal/model"
)

type SignUpResponse struct {
	UserID string `json:"user_id"`
	Token  string `json:"token"`
}

func (s *Server) SignUp(w http.ResponseWriter, r *http.Request) {
	userID := uuid.New()
	err := s.controller.CreateUser(r.Context(), userID)
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to create user: %w", err), http.StatusInternalServerError)
		return
	}

	tokenString, err := s.auth.CreateToken(userID)
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to create jwt token: %w", err), http.StatusInternalServerError)
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
	req := &SignInRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", err), http.StatusBadRequest)
		return
	}
	err = s.controller.GetUser(r.Context(), req.UserID)
	if err == model.ErrNotFound {
		s.ErrorResponse(w, fmt.Errorf("no such user"), http.StatusUnauthorized)
		return
	} else if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to get user: %w", err), http.StatusInternalServerError)
		return
	}
	tokenString, err := s.auth.CreateToken(req.UserID)
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to create jwt token: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, &SignInResponse{Token: tokenString})
}
