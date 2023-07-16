package controller

import (
	"context"

	"github.com/triabokon/gotagv/internal/model"
)

type Storage interface {
	GetUser(ctx context.Context, id string) (string, error)
	InsertUser(ctx context.Context, id string) error

	GetVideo(ctx context.Context, id string) (*model.Video, error)
	ListVideos(ctx context.Context) ([]*model.Video, error)
	InsertVideo(ctx context.Context, video *model.Video) error
	DeleteVideo(ctx context.Context, id string) error

	GetAnnotationWithDuration(ctx context.Context, id string) (*model.Annotation, error)
	ListAnnotations(ctx context.Context) ([]*model.Annotation, error)
	InsertAnnotation(ctx context.Context, a *model.Annotation) error
	UpdateAnnotation(ctx context.Context, id string, p *model.UpdateAnnotationParams) error
	DeleteAnnotation(ctx context.Context, id string) error
}

type Controller struct {
	storage Storage
}

func New(s Storage) *Controller {
	return &Controller{
		storage: s,
	}
}
