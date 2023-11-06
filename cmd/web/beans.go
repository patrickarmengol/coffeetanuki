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
	// parse `id` path parameter
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

	// render template response
	td := newTemplateData()
	td.Bean = bean

	app.render(w, r, http.StatusOK, "beanview.gohtml", "base", td)
}

func (app *application) beanList(w http.ResponseWriter, r *http.Request) {
	beans, err := app.repositories.Beans.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	td := newTemplateData()
	td.Beans = beans

	app.render(w, r, http.StatusOK, "beanlist.gohtml", "base", td)
}

func (app *application) beanCreate(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()
	td.Validator = validator.New()
	td.Bean = &data.Bean{}

	app.render(w, r, http.StatusOK, "beancreate.gohtml", "base", td)
}

func (app *application) beanCreatePost(w http.ResponseWriter, r *http.Request) {
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

	// validate
	v := validator.New()
	bean.Validate(v)

	// invalid form input
	if !v.Valid() {
		td := newTemplateData()
		td.Validator = v
		td.Bean = bean
		app.render(w, r, http.StatusUnprocessableEntity, "beancreate.gohtml", "form", td)
		return
	}

	// try to insert
	err = app.repositories.Beans.Insert(bean)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidRoasterID):
			v.AddFieldError("roaster_id", "roaster id does not exist")
			td := newTemplateData()
			td.Validator = v
			td.Bean = bean
			app.render(w, r, http.StatusUnprocessableEntity, "beancreate.gohtml", "form", td)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// display success message
	td := newTemplateData()
	td.Validator = v
	td.Bean = bean
	td.Result = true
	app.render(w, r, http.StatusOK, "beancreate.gohtml", "form", td)
}
