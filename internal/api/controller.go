package api

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
	"go.uber.org/zap"
)

type Controller struct {
	log     *zap.Logger
	service *service.Service
}

func NewController(log *zap.Logger, service *service.Service) *Controller {
	return &Controller{log, service}
}
