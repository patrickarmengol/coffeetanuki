package main

import (
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/patrickarmengol/coffeetanuki/ui"
)

func (app *application) routes() http.Handler {
	mux := flow.New()

	// static
	fileServer := http.FileServer(http.FS(ui.Files))
	mux.Handle("/static/...", fileServer, http.MethodGet)

	// home
	mux.HandleFunc("/", app.home, http.MethodGet)

	// roaster pages
	mux.HandleFunc("/roasters", app.roasterList, http.MethodGet)
	mux.HandleFunc("/roasters/new", app.roasterCreate, http.MethodGet)
	mux.HandleFunc("/roasters/:id", app.roasterView, http.MethodGet)
	mux.HandleFunc("/roasters/:id/edit", app.roasterEdit, http.MethodGet)
	mux.HandleFunc("/beans/:id", app.beanView, http.MethodGet)

	// roaster htmx
	mux.HandleFunc("/roasters", app.roasterCreatePost, http.MethodPost)
	mux.HandleFunc("/roasters/:id", app.roasterEditPatch, http.MethodPatch)
	mux.HandleFunc("/roasters/:id", app.roasterRemove, http.MethodDelete)

	return mux
}
