package model

import (
	"fmt"
	"time"
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
	ID            string        `json:"id"`
	VideoID       string        `json:"video_id"`
	UserID        string        `json:"user_id"`
	StartTime     time.Duration `json:"start_time"`
	EndTime       time.Duration `json:"end_time"`
	Type          int           `json:"type"`
	Notes         string        `json:"notes"`
	VideoDuration time.Duration `json:"video_duration"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type UpdateAnnotationParams struct {
	StartTime *time.Duration `json:"start_time"`
	EndTime   *time.Duration `json:"end_time"`
	Type      *int           `json:"type"`
	Notes     *string        `json:"notes"`
}

func (p *UpdateAnnotationParams) NoUpdates() bool {
	return p.StartTime == nil &&
		p.EndTime == nil &&
		p.Type == nil &&
		p.Notes == nil
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
	// todo: change to enum
	if p.Type != nil && *p.Type == 0 {
		return fmt.Errorf("type should be specified: %w", ErrInvalidArgument)
	}
	return nil
}
