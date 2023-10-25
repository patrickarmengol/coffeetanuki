package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// overwrite default error response handlers
	// router.NotFound = http.HandlerFunc(app.NotFoundResponse)
	// router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// web routes
	router.HandlerFunc(http.MethodGet, "/", app.home)

	return router
}