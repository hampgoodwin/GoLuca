package api

import (
	"time"

	"github.com/abelgoodwin1988/GoLuca/internal/lucalog"
	"github.com/abelgoodwin1988/GoLuca/internal/setup"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Register api routes and return the router
func Register() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	registerEntryRoute(r)
	registerTransactionRoute(r)
	registerAccountRoutes(r)

	setup.C.Mu.Lock()
	setup.C.Router.Ready = true
	setup.C.Router.Val = r
	setup.C.Mu.Unlock()

	lucalog.Logger.Info("endpoints registered")

	return r
}
