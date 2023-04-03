package service

import (
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type Service struct {
	log        *zap.Logger
	repository *repository.Repository
	publisher  *nats.EncodedConn
}

func NewService(log *zap.Logger, repository *repository.Repository, nec *nats.EncodedConn) *Service {
	return &Service{log, repository, nec}
}
