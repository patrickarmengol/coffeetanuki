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

	// roasters
	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requirePermission("roasters:write"))

		// pages
		mux.HandleFunc("/roasters/new", app.roasterCreate, http.MethodGet)
		mux.HandleFunc("/roasters/:id/edit", app.roasterEdit, http.MethodGet)

		// htmx
		mux.HandleFunc("/roasters", app.roasterCreatePost, http.MethodPost)
		mux.HandleFunc("/roasters/:id", app.roasterEditPatch, http.MethodPatch)
		mux.HandleFunc("/roasters/:id", app.roasterRemove, http.MethodDelete)
	})
	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requirePermission("roasters:read"))

		// pages
		mux.HandleFunc("/roasters", app.roasterList, http.MethodGet)
		mux.HandleFunc("/roasters/:id", app.roasterView, http.MethodGet)
	})

	// beans
	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requirePermission("beans:write"))

		// pages
		mux.HandleFunc("/beans/new", app.beanCreate, http.MethodGet)
		mux.HandleFunc("/beans/:id/edit", app.beanEdit, http.MethodGet)

		// htmx
		mux.HandleFunc("/beans", app.beanCreatePost, http.MethodPost)
		mux.HandleFunc("/beans/:id", app.beanEditPatch, http.MethodPatch)
		mux.HandleFunc("/beans/:id", app.beanRemove, http.MethodDelete)
	})
	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requirePermission("beans:read"))

		// pages
		mux.HandleFunc("/beans", app.beanList, http.MethodGet)
		mux.HandleFunc("/beans/:id", app.beanView, http.MethodGet)
	})

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
