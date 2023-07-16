package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pborman/uuid"

	"github.com/triabokon/gotagv/internal/model"
)

type ListAnnotationsParams struct {
	VideoID string `json:"video_id"`
}

func (c *Controller) ListAnnotations(ctx context.Context, p *ListAnnotationsParams) ([]*model.Annotation, error) {
	if p.VideoID == "" {
		return nil, fmt.Errorf("empty video id: %w", model.ErrInvalidArgument)
	}
	annotations, err := c.storage.ListAnnotations(ctx, p.VideoID)
	if err != nil {
		return nil, fmt.Errorf("failed to list annotations: %w", err)
	}
	return annotations, nil
}

func (c *Controller) CreateAnnotation(ctx context.Context, p *model.CreateAnnotationParams) (string, error) {
	if vErr := p.Validate(); vErr != nil {
		return "", fmt.Errorf("invalid annotation params: %w", vErr)
	}
	video, vErr := c.storage.GetVideo(ctx, p.VideoID)
	if vErr != nil {
		return "", fmt.Errorf("failed to get video: %w", vErr)
	}
	if video.Duration < p.StartTime {
		return "", fmt.Errorf("annotation start time exceeds video duration: %w", model.ErrInvalidArgument)
	}
	if video.Duration < p.EndTime {
		return "", fmt.Errorf("annotation end time exceeds video duration: %w", model.ErrInvalidArgument)
	}

	annotationID := uuid.New()
	annotation := &model.Annotation{
		ID:        annotationID,
		VideoID:   p.VideoID,
		UserID:    p.UserID,
		StartTime: p.StartTime,
		EndTime:   p.EndTime,
		Type:      p.Type,
		Message:   p.Message,
		URL:       p.URL,
		Title:     p.Title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.storage.InsertAnnotation(ctx, annotation); err != nil {
		return "", fmt.Errorf("failed to create annotation: %w", err)
	}
	return annotationID, nil
}

func (c *Controller) UpdateAnnotation(ctx context.Context, id string, p *model.UpdateAnnotationParams) error {
	if p.NoUpdates() {
		return fmt.Errorf("no updates")
	}
	if vErr := p.Validate(); vErr != nil {
		return fmt.Errorf("invalid annotation params: %w", vErr)
	}
	annotation, qErr := c.storage.GetAnnotationWithDuration(ctx, id)
	if qErr != nil {
		return fmt.Errorf("failed to get annotation: %w", qErr)
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
