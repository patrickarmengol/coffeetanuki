package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
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
	td := newTemplateData(r)
	td.Bean = bean

	app.render(w, r, http.StatusOK, "beanview.gohtml", "base", td)
}
