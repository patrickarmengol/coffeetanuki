package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
	"github.com/patrickarmengol/coffeetanuki/internal/validator"
)

// type Bean struct {
// 	ID         int64
// 	Name       string
// 	RoastLevel string
// 	RoasterID  int64
// 	CreatedAt  time.Time
// 	Version    int
// }

type beanForm struct {
	Name       string `form:"name"`
	RoastLevel string `form:"roast_level"`
	RoasterID  int64  `form:"roaster_id"`
}

func (app *application) beanView(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

	// parse id path param
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read bean from db
	bean, err := app.repositories.Beans.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	td.Bean = bean

	// render template response
	app.render(w, r, http.StatusOK, "beanview.gohtml", "base", td)
}

func (app *application) beanList(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

	// read all beans from db
	beans, err := app.repositories.Beans.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	td.Beans = beans

	// render template response
	app.render(w, r, http.StatusOK, "beanlist.gohtml", "base", td)
}

func (app *application) beanCreate(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

	// render form with empty model
	td.Validator = validator.New()
	td.Bean = &data.Bean{}
	app.render(w, r, http.StatusOK, "beancreate.gohtml", "base", td)
}

func (app *application) beanCreatePost(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

	// parse and decode form
	var form beanForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass into model
	bean := &data.Bean{
		Name:       form.Name,
		RoastLevel: form.RoastLevel,
		RoasterID:  form.RoasterID,
	}
	td.Bean = bean

	// validate
	v := validator.New()
	td.Validator = v
	bean.Validate(v)

	// invalid form input
	if !v.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "beancreate.gohtml", "form", td)
		return
	}

	// try to insert
	err = app.repositories.Beans.Insert(bean)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidRoasterID):
			v.AddFieldError("roaster_id", "roaster id does not exist")
			app.render(w, r, http.StatusUnprocessableEntity, "beancreate.gohtml", "form", td)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display success message
	td.Result = true
	app.render(w, r, http.StatusOK, "beancreate.gohtml", "form", td)
}

func (app *application) beanEdit(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read bean from db
	bean, err := app.repositories.Beans.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	td.Bean = bean

	// bogus validator for template
	v := validator.New()
	td.Validator = v

	// render empty form
	app.render(w, r, http.StatusOK, "beanedit.gohtml", "base", td)
}

func (app *application) beanEditPatch(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// read roaster from db
	bean, err := app.repositories.Beans.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	td.Bean = bean

	// decode input form
	var form beanForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass form changes into model
	bean.Name = form.Name
	bean.RoastLevel = form.RoastLevel
	bean.RoasterID = form.RoasterID

	// validate
	v := validator.New()
	td.Validator = v
	bean.Validate(v)

	// case invalid - respond with FieldErrors
	if !v.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "beanedit.gohtml", "form", td)
		return
	}

	// case valid - update roaster
	err = app.repositories.Beans.Update(bean)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidRoasterID):
			v.AddFieldError("roaster_id", "roaster id does not exist")
			app.render(w, r, http.StatusUnprocessableEntity, "beanedit.gohtml", "form", td)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display success
	td.Result = true
	app.render(w, r, http.StatusOK, "beanedit.gohtml", "form", td)
}

func (app *application) beanRemove(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	err = app.repositories.Beans.Delete(id)
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
