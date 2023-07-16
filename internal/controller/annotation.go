package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pborman/uuid"

	"github.com/triabokon/gotagv/internal/model"
)

func (c *Controller) ListAnnotations(ctx context.Context) ([]*model.Annotation, error) {
	annotations, err := c.storage.ListAnnotations(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list annotations: %w", err)
	}
	return annotations, nil
}

type CreateAnnotationParams struct {
	VideoID   string        `json:"video_id"`
	UserID    string        `json:"user_id"`
	StartTime time.Duration `json:"start_time"`
	EndTime   time.Duration `json:"end_time"`
	Type      int           `json:"type"`
	Notes     string        `json:"notes"`
}

func (p *CreateAnnotationParams) Validate() error {
	if p.UserID == "" {
		return fmt.Errorf("empty user id: %w", model.ErrInvalidArgument)
	}
	if p.VideoID == "" {
		return fmt.Errorf("empty video id: %w", model.ErrInvalidArgument)
	}
	if p.StartTime <= 0 {
		return fmt.Errorf("start time should be above 0: %w", model.ErrInvalidArgument)
	}
	if p.EndTime <= 0 {
		return fmt.Errorf("start time should be above 0: %w", model.ErrInvalidArgument)
	}
	if p.EndTime < p.StartTime {
		return fmt.Errorf("start time should be less or equal than end time: %w", model.ErrInvalidArgument)
	}
	// todo: change to enum
	if p.Type == 0 {
		return fmt.Errorf("type should be specified: %w", model.ErrInvalidArgument)
	}
	return nil
}

func (c *Controller) CreateAnnotation(ctx context.Context, p *CreateAnnotationParams) error {
	if vErr := p.Validate(); vErr != nil {
		return fmt.Errorf("invalid annotation params: %w", vErr)
	}
	video, vErr := c.storage.GetVideo(ctx, p.VideoID)
	if vErr != nil {
		return fmt.Errorf("failed to get video: %w", vErr)
	}
	if video.Duration < p.StartTime {
		return fmt.Errorf("annotation start time exceeds video duration: %w", model.ErrInvalidArgument)
	}
	if video.Duration < p.EndTime {
		return fmt.Errorf("annotation end time exceeds video duration: %w", model.ErrInvalidArgument)
	}

	annotation := &model.Annotation{
		ID:        uuid.New(),
		VideoID:   p.VideoID,
		UserID:    p.UserID,
		StartTime: p.StartTime,
		EndTime:   p.EndTime,
		Type:      p.Type,
		Notes:     p.Notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.storage.InsertAnnotation(ctx, annotation); err != nil {
		return fmt.Errorf("failed to create annotation: %w", err)
	}
	return nil
}

func (c *Controller) UpdateAnnotation(ctx context.Context, id string, p *model.UpdateAnnotationParams) error {
	if p.NoUpdates() {
		return fmt.Errorf("no updates")
	}
	if vErr := p.Validate(); vErr != nil {
		return fmt.Errorf("invalid annotation params: %w", vErr)
	}
	annotation, vErr := c.storage.GetAnnotationWithDuration(ctx, id)
	if vErr != nil {
		return fmt.Errorf("failed to get video: %w", vErr)
	}
	if p.StartTime != nil && annotation.VideoDuration < *p.StartTime {
		return fmt.Errorf("annotation start time exceeds video duration: %w", model.ErrInvalidArgument)
	}
	if p.EndTime != nil && annotation.VideoDuration < *p.EndTime {
		return fmt.Errorf("annotation end time exceeds video duration: %w", model.ErrInvalidArgument)
	}

	if err := c.storage.UpdateAnnotation(ctx, id, p); err != nil {
		return fmt.Errorf("failed to update annotation: %w", err)
	}
	return nil
}

func (c *Controller) DeleteAnnotation(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("empty annotation id: %w", model.ErrInvalidArgument)
	}

	if err := c.storage.DeleteAnnotation(ctx, id); err != nil {
		return fmt.Errorf("failed to delete annotation: %w", err)
	}
	return nil
}
