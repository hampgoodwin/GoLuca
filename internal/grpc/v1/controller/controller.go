package controller

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
)

type Controller struct {
	service *service.Service
}

func NewController(s *service.Service) *Controller {
	return &Controller{
		service: s,
	}
}
