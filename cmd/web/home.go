package main

import "net/http"

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)
	app.render(w, r, http.StatusOK, "home.gohtml", "base", td)
}
