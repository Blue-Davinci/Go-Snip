package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	// Create the middleware chain
	globalMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then
	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave).Then
	// Middleware----
	r.Use(globalMiddleware)
	// Home\root path ---
	r.With(dynamicMiddleware).Get("/", app.home)
	// Snippet routes
	r.With(dynamicMiddleware).Mount("/snippet", app.snippetRoutes())
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
	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("./ui/static/")))
	return fileServer
}
