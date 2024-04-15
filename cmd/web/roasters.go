package main

import (
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/errs"
	"github.com/patrickarmengol/coffeetanuki/internal/model"
)

type roasterForm struct {
	Name        string `form:"name"`
	Description string `form:"description"`
	Website     string `form:"website"`
	Location    string `form:"location"`
}

// roaster page
func (app *application) roasterView(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse `id` path parameter
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	// read roaster from db
	roaster, err := app.services.Roasters.Get(r.Context(), id)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.Roaster = roaster

	// render template response
	app.render(w, r, http.StatusOK, "roasterview.gohtml", "base", td)
}

// roaster list page
func (app *application) roasterList(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// decude url query into form
	input := &model.RoasterFilterInput{
		Sort: "id_asc",
	}
	err := app.decodeURLQuery(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid url query format"))
		return
	}
	td.RoasterFilter = input

	// read roasters from service
	roasters, err := app.services.Roasters.Find(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			app.logError(r, err)
			app.render(w, r, http.StatusUnprocessableEntity, "roasterlist.gohtml", "base", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.Roasters = roasters

	// render template response
	app.render(w, r, http.StatusOK, "roasterlist.gohtml", "base", td)
}

// roaster list hx
func (app *application) roasterSearch(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	var input *model.RoasterFilterInput
	err := app.decodeURLQuery(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid url query format"))
		return
	}
	// TODO: figure out how to redirect unprocessable errors to the form

	// read roasters from service
	roasters, err := app.services.Roasters.Find(r.Context(), input)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.Roasters = roasters

	// update client url
	newPath := r.URL.Host + "/roasters?" + r.URL.RawQuery
	w.Header().Add("HX-Push-URL", newPath)

	app.render(w, r, http.StatusOK, "roasterresults.gohtml", "roasterresults", td)
}

// roaster create page
func (app *application) roasterCreate(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.RoasterCreate = &model.RoasterCreateInput{}
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "base", td)
}

// roaster create hx
func (app *application) roasterCreatePost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse and decode form
	var input *model.RoasterCreateInput
	err := app.decodePostForm(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid post form format"))
		return
	}
	td.RoasterCreate = input

	// try to insert
	roaster, err := app.services.Roasters.Create(r.Context(), input)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.Roaster = roaster

	// display success message
	td.Result = true
	app.render(w, r, http.StatusOK, "roastercreate.gohtml", "form", td)
}

// roaster edit page
func (app *application) roasterEdit(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	// read roaster from db
	roaster, err := app.services.Roasters.Get(r.Context(), id)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.RoasterEdit = roaster.ToEditInput()

	// render empty form
	app.render(w, r, http.StatusOK, "roasteredit.gohtml", "base", td)
}

// roaster edit hx
func (app *application) roasterEditPut(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	// decode input form
	input := &model.RoasterEditInput{
		ID: id,
	}
	err = app.decodePostForm(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid post form format"))
		return
	}

	// update roaster
	roaster, err := app.services.Roasters.Update(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			td.RoasterEdit = input // only re-populate input form if validation error
			app.render(w, r, http.StatusUnprocessableEntity, "beanedit.gohtml", "form", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.Roaster = roaster

	// display success
	td.Result = true
	app.render(w, r, http.StatusOK, "roasteredit.gohtml", "form", td)
}

// roaster remove hx
func (app *application) roasterRemove(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	err = app.services.Roasters.Delete(r.Context(), id)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}

	// 200 ok default response
}
