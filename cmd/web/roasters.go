package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
	"github.com/patrickarmengol/coffeetanuki/internal/validator"
)

type roasterForm struct {
	Name                string `form:"name"`
	Description         string `form:"description"`
	Website             string `form:"website"`
	Location            string `form:"location"`
	validator.Validator `form:"-"`
}

func (app *application) roasterView(w http.ResponseWriter, r *http.Request) {
	// parse `id` path parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read roaster from db
	roaster, err := app.models.Roasters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// render template response
	td := newTemplateData(r)
	td.Roaster = roaster

	app.render(w, r, http.StatusOK, "roasterview.gohtml", "base", td)
}

func (app *application) roasterList(w http.ResponseWriter, r *http.Request) {
	roasters, err := app.models.Roasters.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	td := newTemplateData(r)
	td.Roasters = roasters

	app.render(w, r, http.StatusOK, "roasterlist.gohtml", "base", td)
}

func (app *application) roasterCreate(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData(r)
	td.Form = roasterForm{}

	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "base", td)
}

func (app *application) roasterCreatePost(w http.ResponseWriter, r *http.Request) {
	// parse and decode form
	var form roasterForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass into model
	roaster := &data.Roaster{
		Name:        form.Name,
		Description: form.Description,
		Website:     form.Website,
		Location:    form.Location,
	}

	// validate
	roaster.Validate(&form.Validator)

	// case invalid
	if !form.Valid() {
		td := newTemplateData(r)
		td.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "roastercreate.gohtml", "form", td)
		return
	}

	// case valid
	err = app.models.Roasters.Insert(roaster)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// reset form and display success message
	td := newTemplateData(r)
	td.Form = roasterForm{}
	td.Roaster = roaster
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "form", td)
}
