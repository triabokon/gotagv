package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/triabokon/gotagv/internal/auth"
	"github.com/triabokon/gotagv/internal/controller"
	"github.com/triabokon/gotagv/internal/model"
)

type CreateVideoParams struct {
	URL      string `json:"url"`
	Duration string `json:"duration"`
}

func (s *Server) CreateVideo(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		s.ErrorResponse(w, fmt.Errorf("can't access user id"), http.StatusUnauthorized)
		return
	}
	req := &CreateVideoParams{}
	if dErr := json.NewDecoder(r.Body).Decode(req); dErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", dErr), http.StatusBadRequest)
		return
	}
	duration, pErr := parseDuration(req.Duration)
	if pErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse duration: %w", pErr), http.StatusBadRequest)
		return
	}
	err := s.controller.CreateVideo(r.Context(), &controller.CreateVideoParams{
		UserID:   userID,
		URL:      req.URL,
		Duration: duration,
	})
	switch errors.Cause(err) {
	case nil:
	case model.ErrInvalidArgument:
		s.ErrorResponse(w, fmt.Errorf("failed to create video: %w", err), http.StatusBadRequest)
		return
	default:
		s.ErrorResponse(w, fmt.Errorf("failed to create video: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, Response{Code: http.StatusOK, Message: "video created successfully"})
}

func (s *Server) ListVideos(w http.ResponseWriter, r *http.Request) {
	videos, err := s.controller.ListVideos(r.Context())
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to list videos: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, videos)
}

func (s *Server) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	err := s.controller.DeleteVideo(r.Context(), mux.Vars(r)[entityIDKey])
	switch errors.Cause(err) {
	case nil:
	case model.ErrInvalidArgument:
		s.ErrorResponse(w, fmt.Errorf("failed to delete video: %w", err), http.StatusBadRequest)
		return
	default:
		s.ErrorResponse(w, fmt.Errorf("failed to delete video: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, Response{Code: http.StatusOK, Message: "video deleted successfully"})
}
