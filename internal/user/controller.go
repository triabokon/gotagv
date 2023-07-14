package user

import (
	"go.uber.org/zap"
)

type Controller struct {
	logger zap.Logger
}

func New(log zap.Logger) *Controller {
	return &Controller{logger: log}
}
