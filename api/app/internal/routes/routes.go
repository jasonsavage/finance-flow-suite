package routes

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jasonsavage/financeflow/internal/handlers"
	"github.com/jasonsavage/financeflow/internal/repository"
)

func Register(repo repository.DatabaseRepo) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.RequestID)

	// Handlers
	authHandler := handlers.NewAuthHandler(repo)

	// Public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
	})

	// Protected routes
	// r.Group(func(r chi.Router) {
	// 	r.Use(middleware.RequireAuth)

	// 	r.Route("/accounts", func(r chi.Router) {
	// 		r.Get("/", accountHandler.List)
	// 		r.Post("/", accountHandler.Create)
	// 		r.Get("/{id}", accountHandler.Get)
	// 		r.Put("/{id}", accountHandler.Update)
	// 		r.Delete("/{id}", accountHandler.Delete)
	// 	})
	// })

	return r
}
