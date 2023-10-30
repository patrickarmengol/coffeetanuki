package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
	"github.com/patrickarmengol/coffeetanuki/internal/validator"
)

type roasterForm struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Website     string `form:"website"`
	Location    string `form:"location"`
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
	td.Validator = validator.New()
	td.Roaster = &data.Roaster{}

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
	v := validator.New()
	roaster.Validate(v)

	// case invalid
	if !v.Valid() {
		td := newTemplateData(r)
		td.Validator = v
		td.Roaster = roaster
		app.render(w, r, http.StatusUnprocessableEntity, "roastercreate.gohtml", "form", td)
		return
	}

	// case valid
	err = app.models.Roasters.Insert(roaster)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// display success message
	td := newTemplateData(r)
	td.Validator = v
	td.Roaster = roaster
	td.Result = true
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "form", td)
}

func (app *application) roasterEdit(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

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

	v := validator.New()

	td := newTemplateData(r)
	td.Validator = v
	td.Roaster = roaster

	app.render(w, r, http.StatusOK, "roasteredit.gohtml", "base", td)
}

func (app *application) roasterEditPut(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

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

	var form roasterForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass into model
	roaster.Name = form.Name
	roaster.Description = form.Description
	roaster.Website = form.Website
	roaster.Location = form.Location

	v := validator.New()
	roaster.Validate(v)

	if !v.Valid() {
		td := newTemplateData(r)
		td.Validator = v
		td.Roaster = roaster
		app.render(w, r, http.StatusUnprocessableEntity, "roasteredit.gohtml", "form", td)
		return
	}

	err = app.models.Roasters.Update(roaster)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display success message
	td := newTemplateData(r)
	td.Validator = v
	td.Roaster = roaster
	td.Result = true
	app.render(w, r, http.StatusOK, "roasteredit.gohtml", "form", td)
}
