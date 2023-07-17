package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/pborman/uuid"

	"github.com/triabokon/gotagv/internal/model"
)

func (c *Controller) ListVideos(ctx context.Context) ([]*model.Video, error) {
	videos, err := c.storage.ListVideos(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list videos: %w", err)
	}
	return videos, nil
}

type CreateVideoParams struct {
	UserID   string        `json:"user_id"`
	URL      string        `json:"url"`
	Duration time.Duration `json:"duration"`
}

func (p *CreateVideoParams) Validate() error {
	if p.UserID == "" {
		return fmt.Errorf("empty user id: %w", model.ErrInvalidArgument)
	}
	if p.URL == "" {
		return fmt.Errorf("empty url: %w", model.ErrInvalidArgument)
	}
	if p.Duration <= 0 {
		return fmt.Errorf("duration should be above 0: %w", model.ErrInvalidArgument)
	}
	return nil
}

func (c *Controller) CreateVideo(ctx context.Context, p *CreateVideoParams) (string, error) {
	if vErr := p.Validate(); vErr != nil {
		return "", fmt.Errorf("invalid video params: %w", vErr)
	}

	videoID := uuid.New()
	video := &model.Video{
		ID:        videoID,
		UserID:    p.UserID,
		URL:       p.URL,
		Duration:  p.Duration,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := c.storage.InsertVideo(ctx, video); err != nil {
		return "", fmt.Errorf("failed to insert video: %w", err)
	}
	return videoID, nil
}

func (c *Controller) DeleteVideo(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("empty video id: %w", model.ErrInvalidArgument)
	}

	if err := c.storage.DeleteVideo(ctx, id); err != nil {
		return fmt.Errorf("failed to delete video: %w", err)
	}
	return nil
}
