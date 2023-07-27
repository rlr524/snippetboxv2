package main

import "net/http"

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	// StripPrefix is added here as middleware to remove /static from the /static/ routes and
	// hand them over to the NeuteredFileSystem() method to disallow traversing of the static directory.
	mux.Handle("/static/", http.StripPrefix("/static", NeuteredFileSystem(fileServer)))

	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/snippet/view", app.snippetView)
	mux.HandleFunc("/snippet/create", app.snippetCreate)

	return mux
}
