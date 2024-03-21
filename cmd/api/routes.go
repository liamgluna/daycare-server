package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.MethodNotAllowed(app.methodNotAllowedResponse)
	router.NotFound(app.notFoundResponse)

	router.Use(middleware.RealIP)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)
	router.Use(app.rateLimit)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.cfg.allowCORS},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/health", app.healthHandler)

	router.Route("/students", app.loadStudentRoutes)
	router.Route("/faculty", app.loadFacultyRoutes)

	return router
}

func (app *application) loadStudentRoutes(router chi.Router) {
	router.Post("/", app.createStudentWithGuardiansHandler)
	router.Get("/", app.listStudentsHandler)
	router.Get("/{id}", app.showStudentHandler)
	router.Patch("/{id}", app.updateStudentHandler)
	router.Delete("/{id}", app.deleteStudentHandler)
}

func (app *application) loadFacultyRoutes(router chi.Router) {
	router.Post("/", app.createFacultyHandler)
	router.Get("/{id}", app.showFacultyHandler)
	router.Patch("/{id}", app.updateFacultyHandler)
	// router.Delete("/{id}", app.deleteFacultyHandler)
}