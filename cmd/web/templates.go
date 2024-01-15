package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
	"github.com/patrickarmengol/coffeetanuki/internal/validator"
	"github.com/patrickarmengol/coffeetanuki/ui"
)

type templateData struct {
	Bean            *data.Bean
	Beans           []*data.Bean
	Roaster         *data.Roaster
	Roasters        []*data.Roaster
	User            *data.User
	Validator       *validator.Validator
	Result          bool
	IsAuthenticated bool
}

var functions = template.FuncMap{}

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		IsAuthenticated: app.isAuthenticated(r),
	}
}

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

	// get list of partials
	partials, err := fs.Glob(ui.Files, "html/partials/*.gohtml")
	if err != nil {
		return nil, err
	}

	for _, part := range partials {
		name := filepath.Base(part)

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, part)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}

	return cache, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, fileName string, templateName string, data *templateData) {
	// retrieve desired template from cache
	ts, ok := app.templateCache[fileName]
	if !ok {
		err := fmt.Errorf("the template at %s does not exist", fileName)
		app.serverErrorResponse(w, r, err)
		return
	}

	// initialize a buffer, in case of errors executing template
	buf := new(bytes.Buffer)

	// execute template, passing data, and writing to buffer
	// data is currently a pointer to allow for nil; may change later
	err := ts.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// respond with buffer
	w.WriteHeader(status)
	buf.WriteTo(w)
}
