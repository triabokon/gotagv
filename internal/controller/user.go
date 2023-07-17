package controller

import (
	"context"
	"fmt"

	"github.com/triabokon/gotagv/internal/model"
)

func (c *Controller) GetUser(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("empty user id: %w", model.ErrInvalidArgument)
	}

	_, err := c.storage.GetUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	return nil
}

func (c *Controller) CreateUser(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("empty user id: %w", model.ErrInvalidArgument)
	}

	if err := c.storage.InsertUser(ctx, id); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}
