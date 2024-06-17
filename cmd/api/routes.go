package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() *chi.Mux {
	// Initialize a new router
	router := chi.NewRouter()
	//https://github.com/go-chi/cors
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(middleware.Logger)

	// Routes
	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Post("/v1/users", app.createUserHandler)
	router.Get("/v1/users/{id}", app.getUserHandler)
	router.Put("/v1/users/{id}", app.updateUserHandler)
	router.NotFound(app.notFoundResponse)
	return router
}
