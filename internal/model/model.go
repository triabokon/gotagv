package model

import (
	"fmt"
	"time"
)

type AnnotationType string

const (
	UnspecifiedAnnotationType AnnotationType = "unspecified"
	TextAnnotationType        AnnotationType = "text"
	CommentaryAnnotationType  AnnotationType = "commentary"
	LinkAnnotationType        AnnotationType = "link"
	TitleAnnotationType       AnnotationType = "title"
)

type Video struct {
	ID        string        `json:"id"`
	UserID    string        `json:"user_id"`
	URL       string        `json:"url"`
	Duration  time.Duration `json:"duration"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type Annotation struct {
	ID            string         `json:"id"`
	VideoID       string         `json:"video_id"`
	UserID        string         `json:"user_id"`
	StartTime     time.Duration  `json:"start_time"`
	EndTime       time.Duration  `json:"end_time"`
	Type          AnnotationType `json:"type"`
	Message       string         `json:"message,omitempty"`
	URL           string         `json:"url,omitempty"`
	Title         string         `json:"title,omitempty"`
	VideoDuration time.Duration  `json:"video_duration,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type CreateAnnotationParams struct {
	VideoID   string         `json:"video_id"`
	UserID    string         `json:"user_id"`
	StartTime time.Duration  `json:"start_time"`
	EndTime   time.Duration  `json:"end_time"`
	Type      AnnotationType `json:"type"`
	Message   string         `json:"message"`
	URL       string         `json:"url"`
	Title     string         `json:"title"`
}

func (p *CreateAnnotationParams) Validate() error {
	if p.UserID == "" {
		return fmt.Errorf("empty user id: %w", ErrInvalidArgument)
	}
	if p.VideoID == "" {
		return fmt.Errorf("empty video id: %w", ErrInvalidArgument)
	}
	if p.StartTime <= 0 {
		return fmt.Errorf("start time should be above 0: %w", ErrInvalidArgument)
	}
	if p.EndTime <= 0 {
		return fmt.Errorf("start time should be above 0: %w", ErrInvalidArgument)
	}
	if p.EndTime < p.StartTime {
		return fmt.Errorf("start time should be less or equal than end time: %w", ErrInvalidArgument)
	}
	if p.Type == UnspecifiedAnnotationType {
		return fmt.Errorf("invalid type: %w", ErrInvalidArgument)
	}
	return nil
}

type UpdateAnnotationParams struct {
	StartTime *time.Duration  `json:"start_time"`
	EndTime   *time.Duration  `json:"end_time"`
	Type      *AnnotationType `json:"type"`
	Message   *string         `json:"message,omitempty"`
	URL       *string         `json:"url,omitempty"`
	Title     *string         `json:"title,omitempty"`
}

func (p *UpdateAnnotationParams) NoUpdates() bool {
	return p.StartTime == nil &&
		p.EndTime == nil &&
		p.Type == nil &&
		p.Message == nil &&
		p.URL == nil &&
		p.Title == nil
}

func (p *UpdateAnnotationParams) Validate() error {
	if p.StartTime != nil && *p.StartTime == 0 {
		return fmt.Errorf("empty user id: %w", ErrInvalidArgument)
	}
	if p.EndTime != nil && *p.EndTime == 0 {
		return fmt.Errorf("empty video id: %w", ErrInvalidArgument)
	}
	if p.EndTime != nil && p.StartTime != nil && *p.EndTime < *p.StartTime {
		return fmt.Errorf("start time should be less or equal than end time: %w", ErrInvalidArgument)
	}
	if p.Type != nil && *p.Type == UnspecifiedAnnotationType {
		return fmt.Errorf("invalid type: %w", ErrInvalidArgument)
	}
	return nil
}

func ToAnnotationType(t string) AnnotationType {
	switch AnnotationType(t) {
	case TextAnnotationType, CommentaryAnnotationType, LinkAnnotationType, TitleAnnotationType:
		return AnnotationType(t)
	default:
		return UnspecifiedAnnotationType
	}
}
