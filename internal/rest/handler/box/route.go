package boxhandler

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Route("/api/v1/box", func(r chi.Router) {
		r.Get("/connect", h.Connect)
	})
}
