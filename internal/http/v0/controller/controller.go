package controller

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
)

const uuid7Regexp = `^[0-9a-f]{8}(?:\-[0-9a-f]{4}){3}-[0-9a-f]{12}$`

type Controller struct {
	service *service.Service
}

func NewController(service *service.Service) *Controller {
	return &Controller{service}
}
