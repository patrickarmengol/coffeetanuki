package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/patrickarmengol/coffeetanuki/ui"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// load fileserver on embedded static files
	fileServer := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fileServer)

	// overwrite default error response handlers
	// router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// pages
	router.HandlerFunc(http.MethodGet, "/", app.home)

	router.HandlerFunc(http.MethodGet, "/roasters/view/:id", app.roasterView)
	router.HandlerFunc(http.MethodGet, "/roasters/list", app.roasterList)

	return router
}
