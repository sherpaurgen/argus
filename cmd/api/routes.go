package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() *chi.Mux {
	// Initialize a new router
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	// Routes
	router.Get("/v1/healthcheck", app.healthcheckHandler)
	router.Post("/v1/users", app.createUserHandler)
	router.Get("/v1/users/{id}", app.getUserHandler)
	router.NotFound(app.notFoundResponse)
	return router
}
