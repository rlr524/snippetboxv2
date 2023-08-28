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

	r.HandlerFunc(http.MethodGet, "/", app.home)
	r.HandlerFunc(http.MethodGet, "/snippet/view/:id", app.snippetView)
	r.HandlerFunc(http.MethodGet, "/snippet/create", app.snippetCreate)
	r.HandlerFunc(http.MethodPost, "/snippet/create", app.snippetCreatePost)

	// Middleware chain containing the standard middleware which is used for every request
	standard := alice.New(app.recoverPanic, app.logRequests, secureHeaders)

	// Standard middleware chain
	return standard.Then(r)
}
