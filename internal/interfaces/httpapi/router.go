// internal/interfaces/httpapi/router.go
package httpapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *PurchaseHandler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer) // um panic vira 500 em vez de derrubar o servidor

	r.Route("/api/purchases", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/{id}", h.Get)
	})

	return r
}
