package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/alexis/flaggy/internal/sse"
	"github.com/alexis/flaggy/internal/store"
)

// Server holds dependencies for all HTTP handlers.
type Server struct {
	store       store.Store
	broadcaster *sse.Broadcaster
}

// NewRouter creates a Chi router with all routes wired.
func NewRouter(s store.Store, b *sse.Broadcaster) *chi.Mux {
	srv := &Server{store: s, broadcaster: b}

	r := chi.NewRouter()
	r.Use(RequestLogger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {
		// Flags
		r.Post("/flags", srv.CreateFlag)
		r.Get("/flags", srv.ListFlags)
		r.Get("/flags/{key}", srv.GetFlag)
		r.Put("/flags/{key}", srv.UpdateFlag)
		r.Delete("/flags/{key}", srv.DeleteFlag)
		r.Patch("/flags/{key}/toggle", srv.ToggleFlag)

		// Rules
		r.Post("/flags/{key}/rules", srv.CreateRule)
		r.Put("/flags/{key}/rules/{ruleID}", srv.UpdateRule)
		r.Delete("/flags/{key}/rules/{ruleID}", srv.DeleteRule)

		// Evaluate
		r.Post("/evaluate", srv.Evaluate)
		r.Post("/evaluate/batch", srv.EvaluateBatch)

		// SSE Stream
		r.Get("/stream", srv.Stream)
	})

	return r
}
