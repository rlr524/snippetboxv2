package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// The serverError helper writes an error message and stack trace to the errorLog, then sends a
// generic 500 Internal Server Error response to the user. The debug.Stack() function gets a stack
// trace from the current goroutine and appends it to the log message.
func (app *Application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Output(2, trace)

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
