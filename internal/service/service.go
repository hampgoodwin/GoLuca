package service

import (
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"go.uber.org/zap"
)

type Service struct {
	log        *zap.Logger
	repository *repository.Repository
}

func NewService(log *zap.Logger, repository *repository.Repository) *Service {
	return &Service{log, repository}
}
