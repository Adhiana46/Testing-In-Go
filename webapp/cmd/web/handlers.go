package main

import (
	"net/http"
	"path"
	"text/template"
)

var pathToTemplates = "./templates/"

type TemplateData struct {
	IP   string
	Data map[string]any
}

// render page
func (app *application) render(w http.ResponseWriter, r *http.Request, tmpl string, data *TemplateData) error {
	// parse the template from disk
	parsedTemplate, err := template.ParseFiles(path.Join(pathToTemplates, tmpl))
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return err
	}

	// execute the template, passing data if any
	err = parsedTemplate.Execute(w, data)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	_ = app.render(w, r, "home.page.gohtml", &TemplateData{})
}
