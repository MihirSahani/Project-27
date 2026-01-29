package main

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)


func (app *Application) mount() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Timeout(app.config.DefaultContextTimeout))
	router.Use(middleware.Recoverer)
	
	router.Route("/v1", func(router chi.Router) {
		router.Get("/healthcheck", app.healthcheckHandler)
		router.Post("/authenticate", app.authenticationHandler)

		router.Route("/user", func(router chi.Router) {
			router.Post("/", app.createUserHandler)
			router.Group(func(router chi.Router) {
				router.Use(app.authentication)
				router.Delete("/{id}", app.deleteUserHandler)
			})
			
			router.Route("/{id}", func(router chi.Router) {
				router.Use(app.addUserToContext)
				router.Use(app.authentication)
				router.Get("/", app.getUserHandler)
				router.Put("/", app.updateUserHandler)
			})
		})
	})
	return router

}