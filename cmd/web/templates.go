package main

import (
	"bytes"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/errs"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
	"github.com/patrickarmengol/somethingsomethingcoffee/ui"
)

type templateData struct {
	Bean          *model.BeanResponse
	Beans         []*model.BeanResponse
	BeanCreate    *model.BeanCreateInput
	BeanEdit      *model.BeanEditInput
	BeanFilter    *model.BeanFilterInput
	Roaster       *model.RoasterResponse
	Roasters      []*model.RoasterResponse
	RoasterCreate *model.RoasterCreateInput
	RoasterEdit   *model.RoasterEditInput
	RoasterFilter *model.RoasterFilterInput
	User          *model.UserResponse
	UserCreate    *model.UserCreateInput
	UserLogin     *model.UserLoginInput
	// User            *model.User
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
		err := errs.Errorf("the template at %s does not exist", fileName)
		app.errorResponse(w, r, err)
		return
	}

	// initialize a buffer, in case of errors executing template
	buf := new(bytes.Buffer)

	// execute template, passing data, and writing to buffer
	// data is currently a pointer to allow for nil; may change later
	err := ts.ExecuteTemplate(buf, templateName, data)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}

	// respond with buffer
	w.WriteHeader(status)
	buf.WriteTo(w)
}
