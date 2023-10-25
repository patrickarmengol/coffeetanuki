package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/patrickarmengol/coffeetanuki/ui"
)

type templateData struct{}

var functions = template.FuncMap{}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get list of pages
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		// list of template files we want to parse together
		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}

		// register funcmap and parse template files
		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	htmxResponses, err := fs.Glob(ui.Files, "html/htmx/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, response := range htmxResponses {
		name := filepath.Base(response)

		ts, err := template.New(name).ParseFS(ui.Files, response)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *application) renderPage(w http.ResponseWriter, r *http.Request, status int, page string, data *templateData) {
	// retrieve desired template from cache
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverErrorResponse(w, r, err)
		return
	}

	// initialize a buffer
	buf := new(bytes.Buffer)

	// execute with passed data
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(status)
	buf.WriteTo(w)
}
