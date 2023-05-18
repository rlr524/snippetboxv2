package main

import (
	"errors"
	"fmt"
	"github.com/rlr524/snippetboxv2/internal/models"
	"html/template"
	"net/http"
	"strconv"
)

func (app *Application) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	files := []string{
		"./ui/html/base.gohtml",
		"./ui/html/partials/nav.gohtml",
		"./ui/html/pages/home.gohtml",
	}

	// The parameter files... is a variadic parameter meaning it can refer to any number of parameters, in that the
	// slice of file paths in files can be of any length. The ... is the "variadic operator" and
	// works similar to the functionality of the ... spread operator in JavaScript. We can see in the doc
	// that ParseFiles is a variadic function in that it has in its function signature (filenames ...string)
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *Application) SnippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
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
	// Write the snippet data as a plain-text HTTP response body.
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *Application) SnippetCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	// Dummy data
	title := "Oh, snail"
	content := "O snail\nClimb Mount Fuji\nBut slowly, slowly!\n\n- Kobayashi Issa"
	expires := 7

	// Pass the dummy data to the SnippetModel.insert() method, receiving the ID of a new record back
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
