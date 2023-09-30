package handlers

import (
	"github.com/go-chi/chi/v5"

	"github.com/jeronimobarea/go-cqrs/internal/feed"
)

func RegisterRoutes(r *chi.Mux, feedsSvc feed.Service) {
	h := handler{
		feedsSvc: feedsSvc,
	}
	r.Route("/feeds", func(r chi.Router) {
		r.Post("/", h.createFeedHandler)
	})
}
