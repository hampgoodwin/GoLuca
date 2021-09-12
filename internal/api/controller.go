package api

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
	"go.uber.org/zap"
)

const uuidRegexp = "^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}"

type Controller struct {
	log     *zap.Logger
	service *service.Service
}

func NewController(log *zap.Logger, service *service.Service) *Controller {
	return &Controller{log, service}
}
