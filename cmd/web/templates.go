package main

import (
	"github.com/rlr524/snippetboxv2/internal/models"
	"path/filepath"
	"text/template"
	"time"
)

// TemplateData acts as the holding structure for any dynamic data that is passed to the html templates.
type TemplateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
	Form        any
	Flash       string
}

// The humanDate() function returns a formatted string representation of a time.Time object.
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

// Init a template.FuncMap object and store as a global var. This is essentially a string-keyed
// map which acts as a lookup between the names of the custom template functions and the functions themselves.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

// The newTemplateCache() function creates a map for a template cache, loops over all
// file paths, and adds all templates to the cache map.
func newTemplateCache() (map[string]*template.Template, error) {
	// Init a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all file paths that match the app's suffix pattern. This
	// will provide a slice of all the file paths for the application "page" templates.
	pages, err := filepath.Glob("./ui/html/pages/*.go.html")
	if err != nil {
		return nil, err
	}

	// Loop over the file paths
	for _, page := range pages {
		// Extract the file name from the full file path and assign it to the name variable.
		name := filepath.Base(page)

		// The templateFuncMap must be registered with the template set before calling the ParseFiles() method.
		// To do this, use template.New() to create an empty template set, use the template.Funcs() method to register
		// the template,FuncMap, and then parse the file as normal.
		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.go.html")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.go.html")
		if err != nil {
			return nil, err
		}

		// Call ParseFiles() on this template set to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// Add the template set to the map, using the name of the page as the key.
		cache[name] = ts
	}

	return cache, nil
}
