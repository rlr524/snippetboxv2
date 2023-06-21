package main

import "github.com/rlr524/snippetboxv2/internal/models"

// TemplateData acts as the holding structure for any dynamic data that is passed to the html templates.
type TemplateData struct {
	Snippet  *models.Snippet
	Snippets []*models.Snippet
}
