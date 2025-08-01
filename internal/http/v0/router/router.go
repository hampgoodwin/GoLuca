package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

// Register api routes and return the router
func Register(log *zap.Logger, registers ...func(r *chi.Mux)) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httplogger(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	for _, register := range registers {
		register(r)
	}

	return r
}
