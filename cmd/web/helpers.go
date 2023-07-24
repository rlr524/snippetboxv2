package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

// The serverError helper writes an error message and stack trace to the errorLog, then sends a
// generic 500 Internal Server Error response to the user. The debug.Stack() function gets a stack
// trace from the current goroutine and appends it to the log message.
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	e := app.errorLog.Output(2, trace)
	if e != nil {
		return
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// The clientError helper sends a specific code and corresponding description to the user, such as 400
// "Bad Request" responses when there is a problem with a user request. The http.StatusText() function
// generates a human-friendly text representation of a given HTTP status code.
func (app *Application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// The notFound helper is a convenience wrapper around the 404 Not Found client response.
func (app *Application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// The render helper is used to write template data to a buffer then if there are no errors to the
// http.ResponseWriter.
func (app *Application) render(w http.ResponseWriter, status int, page string, data *TemplateData) {
	// Retrieve the appropriate template set from the cache based on the page name. If no entry exists in the cache
	// with the provided name, create a new error and call the serverError() helper method.
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	// Init a new buffer
	buf := new(bytes.Buffer)

	// Write the template to the buffer, instead of straight to the http.ResponseWriter. If there's an error,
	// call the serverError() helper and return.
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Write out the provided HTTP status code.
	w.WriteHeader(status)

	// Write the contents of the buffer to the http.ResponseWriter, in this case passing the http.ResponseWriter
	// to a function that takes in an io.Writer. Explicitly ignore any errors because any error cases have
	// already been handled, and it is guaranteed that there is data in the buffer to write to the http.ResponseWriter.
	_, _ = buf.WriteTo(w)
}

// The newTemplateData helper returns a pointer to a TemplateData struct initialized with the current year.
func (app *Application) newTemplateData(r *http.Request) *TemplateData {
	return &TemplateData{
		CurrentYear: time.Now().Year(),
	}
}
