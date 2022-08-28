package controller

import (
	"github.com/hampgoodwin/GoLuca/internal/service"
	"go.uber.org/zap"
)

type Controller struct {
	log     *zap.Logger
	service service.Service
}
