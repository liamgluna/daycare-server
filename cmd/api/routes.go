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

	// router.Use(middleware.RealIP)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Logger)
	// router.Use(app.rateLimit)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{app.cfg.allowCORS},
		// AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	router.Get("/health", app.healthHandler)

	router.Route("/students", app.loadStudentRoutes)
	router.Route("/faculty", app.loadFacultyRoutes)
	router.Route("/classes", app.loadClassRoutes)
	router.Post("/login", app.loginFacultyHandler)
	router.Post("/logout", app.logoutFacultyHandler)
	return router
}

func (app *application) loadStudentRoutes(router chi.Router) {
	router.Post("/", app.createStudentWithGuardiansHandler)
	router.Get("/", app.listStudentsHandler)
	router.Get("/{id}/guardian", app.showStudentGuardiansHandler)
	router.Patch("/{id}/guardian", app.updateStudentAndGuardianHandler)
	router.Get("/{id}", app.showStudentHandler)
	router.Patch("/{id}", app.updateStudentHandler)
	router.Delete("/{id}", app.deleteStudentHandler)
}

func (app *application) loadFacultyRoutes(router chi.Router) {
	router.Post("/", app.createFacultyHandler)
	router.Get("/{id}", app.showFacultyHandler)
	router.Patch("/profile", app.updateFacultyHandler)
	router.Get("/profile", app.getUserWithTokenHandler)
	// router.Delete("/{id}", app.deleteFacultyHandler)
	
	// get number of classes
	router.Get("/{id}/classes", app.showNumberofClassesByFacultyHandler)
	router.Get("/{id}/classes/students", app.showNumberOStudentsByFacultyHandler)
	router.Get("/{id}/classes/attendance", app.showNumberOfAttendanceTakenByFacultyHandler)
}

func (app *application) loadClassRoutes(router chi.Router) {
	router.Post("/", app.createClassHandler)
	router.Get("/", app.listClassesByFacultyIDHandler)
	router.Get("/{id}", app.showClassHandler)
	router.Patch("/{id}", app.updateClassHandler)
	router.Delete("/{id}", app.deleteClassHandler)

	router.Get("/{classID}/students", app.listClassStudentsHandler)
	router.Post("/{classID}/students", app.createClassStudentHandler)
	router.Delete("/{classID}/students/{studentID}", app.deleteClassStudentHandler)

	router.Post("/{classID}/attendance/{studentID}", app.addStudentAttendance)
	router.Get("/{classID}/attendance", app.getClassAttendance)
}
