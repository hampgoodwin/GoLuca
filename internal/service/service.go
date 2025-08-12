package service

import (
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/nats-io/nats.go"
)

type Service struct {
	repository *repository.Repository
	publisher  *nats.Conn
}

func NewService(repository *repository.Repository, nec *nats.Conn) *Service {
	return &Service{repository, nec}
}
