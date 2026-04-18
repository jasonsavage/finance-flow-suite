package routes

import (
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jasonsavage/financeflow/internal/handlers"
	"github.com/jasonsavage/financeflow/internal/middleware"
	"github.com/jasonsavage/financeflow/internal/repository"
)

func Register(repo repository.DatabaseRepo) *chi.Mux {
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

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

	hcHandler := handlers.NewHealthcheckHandler(repo)
	txHandler := handlers.NewTransactionHandler(repo)
	uHandler := handlers.NewUserHandler(repo)

	// Public routes
	r.Get("/healthcheck", hcHandler.Check)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireAuth)

		r.Get("/user/details", uHandler.GetDetails)
		r.Put("/user/details", uHandler.UpdateDetails)

		r.Route("/transactions", func(r chi.Router) {
			r.Post("/upload", txHandler.UploadTransactions)
			r.Get("/list", txHandler.ListTransactions)
		})

	})

	return r
}
