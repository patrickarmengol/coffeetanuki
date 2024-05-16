package main

import (
	"net/http"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/errs"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
)

// bean page
func (app *application) beanView(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse id path param
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	// read bean from db
	bean, err := app.services.Beans.Get(r.Context(), id)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.Bean = bean

	// render template response
	app.render(w, r, http.StatusOK, "beanview.gohtml", "base", td)
}

// bean list page
func (app *application) beanList(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// decode url query into form
	input := &model.BeanFilterInput{
		Sort: "id_asc",
	}
	err := app.decodeURLQuery(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid url query format"))
		return
	}
	td.BeanFilter = input

	// read beans from db
	beans, err := app.services.Beans.Find(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			app.logError(r, err)
			app.render(w, r, http.StatusUnprocessableEntity, "beanlist.gohtml", "base", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.Beans = beans

	// render template response
	app.render(w, r, http.StatusOK, "beanlist.gohtml", "base", td)
}

// bean list hx
func (app *application) beanSearch(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	var input *model.BeanFilterInput
	err := app.decodeURLQuery(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid url query format"))
		return
	}
	// TODO: figure out how to redirect unprocessable errors to the form

	// read beans from db
	beans, err := app.services.Beans.Find(r.Context(), input)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.Beans = beans

	// update client url
	newPath := r.URL.Host + "/beans?" + r.URL.RawQuery
	w.Header().Add("HX-Push-URL", newPath)

	// render template response
	app.render(w, r, http.StatusOK, "beanresults.gohtml", "beanresults", td)
}

// bean create page
func (app *application) beanCreate(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.BeanCreate = &model.BeanCreateInput{}
	app.render(w, r, http.StatusOK, "beancreate.gohtml", "base", td)
}

// bean create hx
func (app *application) beanCreatePost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse and decode form
	var input *model.BeanCreateInput
	err := app.decodePostForm(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid post form format"))
		return
	}
	td.BeanCreate = input

	// try to insert
	bean, err := app.services.Beans.Create(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			app.render(w, r, http.StatusUnprocessableEntity, "beancreate.gohtml", "form", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.Bean = bean

	// display success message
	td.Result = true
	app.render(w, r, http.StatusOK, "beancreate.gohtml", "form", td)
}

// bean edit page
func (app *application) beanEdit(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	// read bean from db
	bean, err := app.services.Beans.Get(r.Context(), id)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}
	td.BeanEdit = bean.ToEditInput()

	// render empty form
	app.render(w, r, http.StatusOK, "beanedit.gohtml", "base", td)
}

// bean edit hx
func (app *application) beanEditPut(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// read id from path
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	// decode input form
	input := &model.BeanEditInput{
		ID: id,
	}
	err = app.decodePostForm(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid post form format"))
		return
	}

	// update roaster
	bean, err := app.services.Beans.Update(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			td.BeanEdit = input // only re-populate input form if validation error
			app.render(w, r, http.StatusUnprocessableEntity, "beanedit.gohtml", "form", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.Bean = bean

	// display success
	td.Result = true
	app.render(w, r, http.StatusOK, "beanedit.gohtml", "form", td)
}

// bean remove hx
func (app *application) beanRemove(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid id format"))
		return
	}

	err = app.services.Beans.Delete(r.Context(), id)
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}

	// 200 ok default response
}
