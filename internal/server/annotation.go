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

type CreateAnnotationRequest struct {
	VideoID   string `json:"video_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Type      string `json:"type"`
	Message   string `json:"message"`
	URL       string `json:"url"`
	Title     string `json:"title"`
}

type CreateAnnotationResponse struct {
	AnnotationID string `json:"annotation_id"`
}

func toCreateAnnotationParams(r *CreateAnnotationRequest, userID string) (*model.CreateAnnotationParams, error) {
	startTime, pErr := parseDuration(r.StartTime)
	if pErr != nil {
		return nil, fmt.Errorf("failed to parse start time: %w", pErr)
	}
	endTime, pErr := parseDuration(r.EndTime)
	if pErr != nil {
		return nil, fmt.Errorf("failed to parse end time: %w", pErr)
	}
	p := &model.CreateAnnotationParams{
		VideoID:   r.VideoID,
		UserID:    userID,
		StartTime: startTime,
		EndTime:   endTime,
		Type:      model.ToAnnotationType(r.Type),
		Message:   r.Message,
		URL:       r.URL,
		Title:     r.Title,
	}
	return p, nil
}

func (s *Server) CreateAnnotation(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(auth.UserIDKey).(string)
	if !ok {
		s.ErrorResponse(w, fmt.Errorf("no user id"), http.StatusUnauthorized)
		return
	}
	req := &CreateAnnotationRequest{}
	if dErr := json.NewDecoder(r.Body).Decode(req); dErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", dErr), http.StatusBadRequest)
		return
	}
	p, pErr := toCreateAnnotationParams(req, userID)
	if pErr != nil {
		s.ErrorResponse(w, pErr, http.StatusBadRequest)
		return
	}
	annotationID, err := s.controller.CreateAnnotation(r.Context(), p)
	if errors.Is(err, model.ErrInvalidArgument) || errors.Is(err, model.ErrAlreadyExists) {
		s.ErrorResponse(w, fmt.Errorf("failed to create annotation: %w", err), http.StatusBadRequest)
		return
	}
	if errors.Is(err, model.ErrNotFound) {
		s.ErrorResponse(w, fmt.Errorf("failed to create annotation: %w", err), http.StatusNotFound)
		return
	}
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to create annotation: %w", err), http.StatusInternalServerError)
		return
	}
	s.SuccessResponse(w, CreateAnnotationResponse{AnnotationID: annotationID})
}

type UpdateAnnotationRequest struct {
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
	Type      *string `json:"type,omitempty"`
	Message   *string `json:"message,omitempty"`
	URL       *string `json:"url,omitempty"`
	Title     *string `json:"title,omitempty"`
}

func toUpdateAnnotationParams(r *UpdateAnnotationRequest) (*model.UpdateAnnotationParams, error) {
	var aType *model.AnnotationType
	if r.Type != nil {
		t := model.ToAnnotationType(*r.Type)
		aType = &t
	}
	p := &model.UpdateAnnotationParams{
		Type: aType, Message: r.Message, URL: r.URL, Title: r.Title,
	}
	if r.StartTime != nil {
		startTime, pErr := parseDuration(*r.StartTime)
		if pErr != nil {
			return nil, fmt.Errorf("failed to parse start time: %w", pErr)
		}
		p.StartTime = &startTime
	}
	if r.EndTime != nil {
		endTime, pErr := parseDuration(*r.EndTime)
		if pErr != nil {
			return nil, fmt.Errorf("failed to parse start time: %w", pErr)
		}
		p.EndTime = &endTime
	}
	return p, nil
}

func (s *Server) UpdateAnnotation(w http.ResponseWriter, r *http.Request) {
	req := &UpdateAnnotationRequest{}
	if dErr := json.NewDecoder(r.Body).Decode(req); dErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", dErr), http.StatusBadRequest)
		return
	}
	p, pErr := toUpdateAnnotationParams(req)
	if pErr != nil {
		s.ErrorResponse(w, pErr, http.StatusBadRequest)
		return
	}
	err := s.controller.UpdateAnnotation(r.Context(), mux.Vars(r)[entityIDKey], p)
	if errors.Is(err, model.ErrInvalidArgument) || errors.Is(err, model.ErrAlreadyExists) {
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusBadRequest)
		return
	}
	if errors.Is(err, model.ErrNotFound) {
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusNotFound)
		return
	}
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusInternalServerError)
		return
	}
	s.SuccessResponse(w, Response{Message: "annotation updated successfully"})
}

type ListAnnotationsResponse struct {
	Annotations []*model.Annotation `json:"annotations"`
}

func (s *Server) ListAnnotations(w http.ResponseWriter, r *http.Request) {
	req := &controller.ListAnnotationsParams{}
	if dErr := json.NewDecoder(r.Body).Decode(req); dErr != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to parse request: %w", dErr), http.StatusBadRequest)
		return
	}
	annotations, err := s.controller.ListAnnotations(r.Context(), req)
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to list annotations: %w", err), http.StatusInternalServerError)
		return
	}
	s.SuccessResponse(w, &ListAnnotationsResponse{Annotations: annotations})
}

func (s *Server) DeleteAnnotation(w http.ResponseWriter, r *http.Request) {
	err := s.controller.DeleteAnnotation(r.Context(), mux.Vars(r)[entityIDKey])
	if errors.Is(err, model.ErrInvalidArgument) {
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusBadRequest)
		return
	}
	if err != nil {
		s.ErrorResponse(w, fmt.Errorf("failed to update annotation: %w", err), http.StatusInternalServerError)
		return
	}
	s.SuccessResponse(w, Response{Message: "annotation deleted successfully"})
}
