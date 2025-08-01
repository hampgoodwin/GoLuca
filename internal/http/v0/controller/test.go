package controller

import (
	"net/http"

	"github.com/hampgoodwin/GoLuca/internal/http/v0/router"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func newTestHTTPHandler(
	log *zap.Logger,
	db *pgxpool.Pool,
) http.Handler {
	r := repository.NewRepository(db)
	nc, _ := nats.Connect(nats.DefaultURL)
	s := service.NewService(log, r, nc)
	c := NewController(log, s)
	return router.Register(
		log,
		c.RegisterAccountRoutes,
		c.RegisterTransactionRoutes,
	)
}
