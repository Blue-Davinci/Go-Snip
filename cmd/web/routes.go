package main

import (
	"net/http"

	"github.com/blue-davinci/gosnip/ui"
	"github.com/go-chi/chi/v5"
	"github.com/justinas/alice"
)

// routes() returns a http.Handler which defines the routes for the application.
func (app *application) routes() http.Handler {
	r := chi.NewRouter()
	// Create the middleware chain
	globalMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders).Then
	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate).Then
	// Add Global Middleware----
	r.Use(globalMiddleware)
	// Home\root path ---
	r.With(dynamicMiddleware).Get("/", app.home)
	// Snippet routes
	r.With(dynamicMiddleware).Mount("/snippet", app.snippetRoutes())
	// User routes
	r.With(dynamicMiddleware).Mount("/user", app.userRoutes())
	// Static file handler to the serve fs static files
	r.Handle("/static/*", staticFileHandlerRoute())
	// Http error Handlers
	app.setCustomHttpErrorHandlers(r)

	return r
}

// snippetRoutes() returns a chi.Router which defines the snippet routes
// for the application.
func (app *application) snippetRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/view/{id}", app.snippetView)
	// create
	r.With(app.requireAuthentication).Post("/create", app.snippetCreatePost)
	r.With(app.requireAuthentication).Get("/create", app.snippetCreate)
	return r
}

// userRoutes() returns a chi.Router which defines the user routes
// for the application.
func (app *application) userRoutes() chi.Router {
	r := chi.NewRouter()
	r.Get("/signup", app.userSignup)
	r.Post("/signup", app.userSignupPost)
	r.Get("/login", app.userLogin)
	r.Post("/login", app.userLoginPost)
	//
	r.With(app.requireAuthentication).Post("/logout", app.userLogoutPost)
	return r
}

// setCustomHttpErrorHandlers() sets custom http error handlers
// for 404 and 405 errors
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

// staticFileHandlerRoute() returns a http.Handler which serves the static files
// from the embedded filesystem.
func staticFileHandlerRoute() http.Handler {
	// Take the ui.Files embedded filesystem and convert it to a http.FS type so
	// that it satisfies the http.FileSystem interface. We then pass that to the
	// http.FileServer() function to create the file server handler.
	fileServer := http.FileServer(http.FS(ui.Files))
	return fileServer
}
