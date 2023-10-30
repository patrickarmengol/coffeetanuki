package main

import (
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) errorResponse(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// 500 - Internal Server Error
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	app.errorResponse(w, http.StatusInternalServerError)
}

// 400 - Bad Request
func (app *application) badRequestResponse(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusBadRequest)
}

// 404 - Not Found
func (app *application) notFoundResponse(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusNotFound)
}

// 409 - Edit Conflict
func (app *application) editConflictResponse(w http.ResponseWriter) {
	app.errorResponse(w, http.StatusConflict)
}
