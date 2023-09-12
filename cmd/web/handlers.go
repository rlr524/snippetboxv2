package main

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rlr524/snippetboxv2/internal/models"
	"net/http"
	"strconv"
)

func (app *Application) home(w http.ResponseWriter, r *http.Request) {
	snippets, err := app.snippets.GetLatest()
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Use the newTemplateData() helper to get a TemplateData struct containing the "default" data and
	// add the snippets slice to it.
	data := app.newTemplateData(r)
	data.Snippets = snippets

	app.render(w, http.StatusOK, "home.gohtml", data)
}

func (app *Application) snippetView(w http.ResponseWriter, r *http.Request) {
	// When httprouter is parsing a request, the values of any named parameters will be stored in the request context.
	params := httprouter.ParamsFromContext(r.Context())

	// Then use the ByName() method and get the value of the parameter named "id" from the slice and validate.
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	// Use the SnippetModel's Get method to retrieve the data for a specific record based on its ID. If no
	// matching record is found, return a 404 Not Found response.
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, http.StatusOK, "view.gohtml", data)
}

func (app *Application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	app.render(w, http.StatusOK, "create.gohtml", data)
}

func (app *Application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	// Call r.ParseForm(), which adds any data in the POST request bodies to the r.PostForm map.
	// This works the same way for put and patch requests. If there are any errors, use the app.ClientError()
	// helper to send a 400 Bad Request to the user.
	r.Body = http.MaxBytesReader(w, r.Body, 4096)

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Use the r.PostForm.Get() method to retrieve the title and content from the r.PostForm map.
	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")

	// The PostForm.Get() method always returns the form data as a string. However, we expect the "expires" value
	// to be a number and want to represent it in code as an int. So manually convert the form data to an int
	// using strconv.Atoi() and send a 400 if the conversion fails.
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}
