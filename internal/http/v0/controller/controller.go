package controller

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
	"go.uber.org/zap"
)

const ksuidRegexp = "^[a-zA-Z0-9]{27}$"

type Controller struct {
	log     *zap.Logger
	service *service.Service
}

func NewController(log *zap.Logger, service *service.Service) *Controller {
	return &Controller{log, service}
}
