package main

import (
	"github.com/justinas/alice"
	"net/http"
)

// The Routes method instantiates a new ServeMux from the net/http package, sets the static file server directory,
// invokes the NeuteredFileSystem function, and handles all routes, returning a http.Handler.
func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// StripPrefix is added here as middleware to remove /static from the /static/ routes and
	// hand them over to the NeuteredFileSystem() method to disallow traversing of the static directory.
	mux.Handle("/static/", http.StripPrefix("/static", NeuteredFileSystem(fileServer)))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	// Middleware chain containing the standard middleware which is used for every request
	standard := alice.New(app.recoverPanic, app.logRequests, secureHeaders)

	// Standard middleware chain
	return standard.Then(mux)
}
