package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	// Middleware----
	app.loadMiddleware(r)
	// Home\root path ---
	r.Get("/", app.home)
	// Snippet routes
	r.Mount("/snippet", app.snippetRoutes())
	// Static file handler to the router
	r.Handle("/static/*", staticFileHandlerRoute())
	// Http error Handlers
	app.setCustomHttpErrorHandlers(r)

	return r
}
func (app *application) snippetRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/view/{id}", app.snippetView)
	// create
	r.Post("/create", app.snippetCreatePost)
	r.Get("/create", app.snippetCreate)
	return r
}

func (app *application) setCustomHttpErrorHandlers(r *chi.Mux) {
	// Set a custom 404 handler
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})
	// Set a custom 405 handler
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		app.clientError(w, http.StatusMethodNotAllowed)
	})
}
func staticFileHandlerRoute() http.Handler {
	// Create a file server which serves files out of the "./ui/static" directory.
	// Note that the path given to the http.Dir function is relative to the project
	// directory root.
	fileServer := http.StripPrefix("/static/*filepath", http.FileServer(http.Dir("./ui/static/")))
	return fileServer
}

func (app *application) loadMiddleware(r *chi.Mux) {
	// Add panic recovery middleware
	r.Use(app.recoverPanic)
	// Add the logger middleware
	r.Use(app.logRequest)
	// Add the secureHeaders middleware
	r.Use(secureHeaders)
}
