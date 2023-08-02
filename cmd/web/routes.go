package main

import "net/http"

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

	// Wrap the mux in the secureHeaders middleware as the "next" parameter;
	// wrapped in the logRequest middleware;
	// wrapped in the recoverPanic middleware.
	return app.recoverPanic(app.logRequests(secureHeaders(mux)))
}
