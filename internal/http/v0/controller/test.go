package controller

import (
	"net/http"

	"github.com/hampgoodwin/GoLuca/internal/http/v0/router"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/encoders/protobuf"
	"go.uber.org/zap"
)

func newTestHTTPHandler(
	log *zap.Logger,
	db *pgxpool.Pool,
) http.Handler {
	r := repository.NewRepository(db)
	nc, _ := nats.Connect(nats.DefaultURL)
	nec, _ := nats.NewEncodedConn(nc, protobuf.PROTOBUF_ENCODER)
	s := service.NewService(log, r, nec)
	c := NewController(log, s)
	return router.Register(
		log,
		c.RegisterAccountRoutes,
		c.RegisterTransactionRoutes,
	)
}
