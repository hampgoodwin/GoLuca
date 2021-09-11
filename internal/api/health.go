package api

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (c *Controller) RegisterHealthRoutes(r *chi.Mux) {
	r.Get("/health", c.health)
}

func (c *Controller) health(w http.ResponseWriter, r *http.Request) {
	c.log.Debug("health check")
	c.respond(w, struct{}{}, http.StatusOK)
}
