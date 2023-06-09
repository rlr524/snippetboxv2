package main

import "net/http"

func (app *Application) Routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	//fileServer := http.FileServer(http.Dir(app.cfg.staticDir))
	mux.Handle("/static/", http.StripPrefix("/static", NeuteredFileSystem(fileServer)))

	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/snippet/view", app.SnippetView)
	mux.HandleFunc("/snippet/create", app.SnippetCreate)

	return mux
}
