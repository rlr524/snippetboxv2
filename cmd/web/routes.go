package main

import (
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

// The Routes method instantiates a new ServeMux from the net/http package, sets the static file server directory,
// invokes the NeuteredFileSystem function, and handles all routes, returning a http.Handler.
func (app *Application) Routes() http.Handler {
	r := httprouter.New()

	// The NotFound handler takes in a closure that returns the custom notFound message
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// StripPrefix is added here as middleware to remove /static from the /static/ routes and
	// hand them over to the NeuteredFileSystem() method to disallow traversing of the static directory.
	r.Handler(http.MethodGet,
		"/static/*filepath",
		http.StripPrefix("/static",
			NeuteredFileSystem(fileServer)))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	r.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	r.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(app.snippetView))
	r.Handler(http.MethodGet, "/snippet/create", dynamic.ThenFunc(app.snippetCreate))
	r.Handler(http.MethodPost, "/snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	// Middleware chain containing the standard middleware which is used for every request
	standard := alice.New(app.recoverPanic, app.logRequests, secureHeaders)

	// Standard middleware chain
	return standard.Then(r)
}
