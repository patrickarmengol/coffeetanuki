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

	// use session management for dynamic content
	mux.Use(app.sessionManager.LoadAndSave)

	// use authentication
	mux.Use(app.authenticate)

	// home
	mux.HandleFunc("/", app.home, http.MethodGet)

	// roaster pages
	mux.HandleFunc("/roasters", app.roasterList, http.MethodGet)
	mux.HandleFunc("/roasters/new", app.roasterCreate, http.MethodGet)
	mux.HandleFunc("/roasters/:id", app.roasterView, http.MethodGet)
	mux.HandleFunc("/roasters/:id/edit", app.roasterEdit, http.MethodGet)

	// roaster htmx
	mux.HandleFunc("/roasters", app.roasterCreatePost, http.MethodPost)
	mux.HandleFunc("/roasters/:id", app.roasterEditPatch, http.MethodPatch)
	mux.HandleFunc("/roasters/:id", app.roasterRemove, http.MethodDelete)

	// bean pages
	mux.HandleFunc("/beans", app.beanList, http.MethodGet)
	mux.HandleFunc("/beans/new", app.beanCreate, http.MethodGet)
	mux.HandleFunc("/beans/:id", app.beanView, http.MethodGet)
	mux.HandleFunc("/beans/:id/edit", app.beanEdit, http.MethodGet)

	// bean htmx
	mux.HandleFunc("/beans", app.beanCreatePost, http.MethodPost)
	mux.HandleFunc("/beans/:id", app.beanEditPatch, http.MethodPatch)
	mux.HandleFunc("/beans/:id", app.beanRemove, http.MethodDelete)

	// user pages
	mux.HandleFunc("/user/signup", app.userSignup, http.MethodGet)
	mux.HandleFunc("/user/login", app.userLogin, http.MethodGet)
	mux.HandleFunc("/account", app.userAccountView, http.MethodGet)

	// user htmx
	mux.HandleFunc("/user/signup", app.userSignupPost, http.MethodPost)
	mux.HandleFunc("/user/login", app.userLoginPost, http.MethodPost)
	mux.HandleFunc("/user/logout", app.userLogoutPost, http.MethodPost)

	return mux
}
