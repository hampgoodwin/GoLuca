package controller

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
)

const ksuidRegexp = "^[a-zA-Z0-9]{27}$"

type Controller struct {
	service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{service}
}
