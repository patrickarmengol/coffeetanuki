package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
)

func (app *application) roasterView(w http.ResponseWriter, r *http.Request) {
	// parse `id` path parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// read roaster from db
	roaster, err := app.models.Roasters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// render template response
	td := newTemplateData(r)
	td.Roaster = roaster

	app.render(w, r, http.StatusOK, "roasterview.gohtml", "base", &td)
}

func (app *application) roasterList(w http.ResponseWriter, r *http.Request) {
	roasters, err := app.models.Roasters.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	td := newTemplateData(r)
	td.Roasters = roasters

	app.render(w, r, http.StatusOK, "roasterlist.gohtml", "base", &td)
}

func (app *application) roasterCreate(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "base", nil)
}
