package main

import (
	"github.com/rlr524/snippetboxv2/internal/models"
	"path/filepath"
	"text/template"
)

// TemplateData acts as the holding structure for any dynamic data that is passed to the html templates.
type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// Init a new map to act as the cache
	cache := map[string]*template.Template{}

	// Use the filepath.Glob() function to get a slice of all file paths that match the app's suffix pattern. This
	// will provide a slice of all the file paths for the application "page" templates.
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	// Loop over the file paths
	for _, page := range pages {
		// Extract the file name from the full file path and assign it to the name variable.
		name := filepath.Base(page)

		// Parse the base template file into the template set.
		ts, err := template.ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}

		// Call ParseGlob() on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.gohtml")
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
