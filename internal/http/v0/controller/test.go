package controller

import (
	"net/http"

	"github.com/hampgoodwin/GoLuca/internal/http/router"
	"github.com/hampgoodwin/GoLuca/internal/repository"
	"github.com/hampgoodwin/GoLuca/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

func newTestHTTPHandler(
	log *zap.Logger,
	db *pgxpool.Pool,
) http.Handler {
	r := repository.NewRepository(db)
	s := service.NewService(log, r)
	c := NewController(log, s)
	return router.Register(
		log,
		c.RegisterAccountRoutes,
		c.RegisterTransactionRoutes,
	)
}
