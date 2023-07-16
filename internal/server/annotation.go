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

type CreateAnnotationParams struct {
	VideoID   string `json:"video_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      int    `json:"type"`
	Notes     string `json:"notes"`
}

func (s *Server) CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		s.ErrorResponse(w, fmt.Errorf("no user id"), http.StatusUnauthorized)
		return
	}
	req := &CreateAnnotationParams{}
	if dErr := json.NewDecoder(r.Body).Decode(req); dErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", dErr), http.StatusBadRequest)
		return
	}
	startTime, pErr := parseDuration(req.StartTime)
	if pErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse start time: %w", pErr), http.StatusBadRequest)
		return
	}
	endTime, pErr := parseDuration(req.EndTime)
	if pErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse start time: %w", pErr), http.StatusBadRequest)
		return
	}
	err := s.controller.CreateAnnotation(r.Context(), &controller.CreateAnnotationParams{
		VideoID:   req.VideoID,
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
		Type:      req.Type,
		Notes:     req.Notes,
	})
	switch errors.Cause(err) {
	case nil:
	case model.ErrInvalidArgument:
		s.ErrorResponse(w, fmt.Errorf("failed to create annotation: %w", err), http.StatusBadRequest)
		return
	default:
		s.ErrorResponse(w, fmt.Errorf("failed to create annotation: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, Response{Code: http.StatusOK, Message: "annotation created successfully"})
}

type UpdateAnnotationParams struct {
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Type      *int    `json:"type"`
	Notes     *string `json:"notes"`
}

func (s *Server) UpdateAnnotation(w http.ResponseWriter, r *http.Request) {
	req := &UpdateAnnotationParams{}
	if dErr := json.NewDecoder(r.Body).Decode(req); dErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", dErr), http.StatusBadRequest)
		return
	}
	p := &model.UpdateAnnotationParams{Type: req.Type, Notes: req.Notes}
	if req.StartTime != nil {
		startTime, pErr := parseDuration(*req.StartTime)
		if pErr != nil {
			s.ErrorResponse(w, fmt.Errorf("failed to parse start time: %w", pErr), http.StatusBadRequest)
			return
		}
		p.StartTime = &startTime
	}
	if req.EndTime != nil {
		endTime, pErr := parseDuration(*req.EndTime)
		if pErr != nil {
			s.ErrorResponse(w, fmt.Errorf("failed to parse end time: %w", pErr), http.StatusBadRequest)
			return
		}
		p.EndTime = &endTime
	}
	err := s.controller.UpdateAnnotation(r.Context(), mux.Vars(r)[entityIDKey], p)
	switch errors.Cause(err) {
	case nil:
	case model.ErrInvalidArgument, model.ErrAlreadyExists:
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusBadRequest)
		return
	case model.ErrNotFound:
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusNotFound)
		return
	default:
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, Response{Code: http.StatusOK, Message: "annotation updated successfully"})
}

func (s *Server) ListAnnotations(w http.ResponseWriter, r *http.Request) {
	annotations, err := s.controller.ListAnnotations(r.Context())
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to list annotations: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, annotations)
}

func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	err := s.controller.DeleteAnnotation(r.Context(), mux.Vars(r)[entityIDKey])
	switch errors.Cause(err) {
	case nil:
	case model.ErrInvalidArgument:
		s.ErrorResponse(w, fmt.Errorf("failed to delete annotation: %w", err), http.StatusBadRequest)
		return
	default:
		s.ErrorResponse(w, fmt.Errorf("failed to delete annotation: %w", err), http.StatusInternalServerError)
		return
	}
	s.JSONResponse(w, Response{Code: http.StatusOK, Message: "annotation deleted successfully"})
}
