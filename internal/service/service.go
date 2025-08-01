package service

import (
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Service struct {
	log        *zap.Logger
	repository *repository.Repository
	publisher  *nats.Conn
}

func NewService(log *zap.Logger, repository *repository.Repository, nec *nats.Conn) *Service {
	return &Service{log, repository, nec}
}
