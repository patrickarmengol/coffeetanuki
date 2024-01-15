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
	td := app.newTemplateData(r)

	// parse `id` path parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read roaster from db
	roaster, err := app.repositories.Roasters.GetFull(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	td.Roaster = roaster

	// render template response
	app.render(w, r, http.StatusOK, "roasterview.gohtml", "base", td)
}

func (app *application) roasterList(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read roasters from db
	roasters, err := app.repositories.Roasters.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	td.Roasters = roasters

	// render template response
	app.render(w, r, http.StatusOK, "roasterlist.gohtml", "base", td)
}

func (app *application) roasterSearch(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	var input struct {
		SearchTerm string `form:"search"`
		PageNum    int    `form:"page"`
		PageSize   int    `form:"page_size"`
		Sort       string `form:"sort"`
	}

	// TODO: search input validation

	err := app.decodeURLQuery(r, &input)
	if err != nil {
		app.badRequestResponse(w)
		return
	}
	roasters, err := app.repositories.Roasters.Search(input.SearchTerm, input.PageNum, input.PageSize, input.Sort)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	td.Roasters = roasters

	app.render(w, r, http.StatusOK, "roasterresults.gohtml", "roasterresults", td)
}

func (app *application) roasterCreate(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.Validator = validator.New()
	td.Roaster = &data.Roaster{}
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "base", td)
}

func (app *application) roasterCreatePost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

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
	td.Roaster = roaster

	// validate
	v := validator.New()
	td.Validator = v
	roaster.Validate(v)

	// case invalid - respond with FieldErrors
	if !v.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "roastercreate.gohtml", "form", td)
		return
	}

	// case valid
	err = app.repositories.Roasters.Insert(roaster)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// display success message
	td.Result = true
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "form", td)
}

func (app *application) roasterEdit(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read roaster from db
	roaster, err := app.repositories.Roasters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	td.Roaster = roaster

	// bogus validator for template
	v := validator.New()
	td.Validator = v

	// render empty form
	app.render(w, r, http.StatusOK, "roasteredit.gohtml", "base", td)
}

func (app *application) roasterEditPatch(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read roaster from db
	roaster, err := app.repositories.Roasters.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	td.Roaster = roaster

	// decode input form
	var form roasterForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass form changes into model
	roaster.Name = form.Name
	roaster.Description = form.Description
	roaster.Website = form.Website
	roaster.Location = form.Location

	// validate
	v := validator.New()
	td.Validator = v
	roaster.Validate(v)

	// case invalid - respond with FieldErrors
	if !v.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "roasteredit.gohtml", "form", td)
		return
	}

	// case valid - update roaster
	err = app.repositories.Roasters.Update(roaster)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display success
	td.Result = true
	app.render(w, r, http.StatusOK, "roasteredit.gohtml", "form", td)
}

func (app *application) roasterRemove(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	err = app.repositories.Roasters.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// 200 ok default response
}
