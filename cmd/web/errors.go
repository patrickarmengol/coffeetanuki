package main

import (
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/errs"
)

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, err error) {
	code, message := errs.ErrorCode(err), errs.ErrorMessage(err)

	if code == errs.ERRINTERNAL {
		app.logError(r, err)
	}

	status := errorStatus(code)

	// TODO: maybe do a template response?
	// figure out if i want to expose the error message to user
	// http.Error(w, http.StatusText(status), status)
	// http.Error(w, err.Error(), status)
	http.Error(w, message, status)
}

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

var codetostatus = map[string]int{
	errs.ERRCONFLICT:       http.StatusConflict,
	errs.ERRBAD:            http.StatusBadRequest,
	errs.ERRNOTFOUND:       http.StatusNotFound,
	errs.ERRNOTIMPLEMENTED: http.StatusNotImplemented,
	errs.ERRNOTAUTHORIZED:  http.StatusUnauthorized,
	errs.ERRINTERNAL:       http.StatusInternalServerError,
}

func errorStatus(errCode string) int {
	if v, ok := codetostatus[errCode]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// func (app *application) logError(r *http.Request, err error) {
// 	var (
// 		method = r.Method
// 		uri    = r.URL.RequestURI()
// 	)
//
// 	app.logger.Error(err.Error(), "method", method, "uri", uri)
// }
//
// TODO: expand these responses to full pages
//
// func (app *application) errorResponse(w http.ResponseWriter, status int) {
// 	http.Error(w, http.StatusText(status), status)
// }
//
// // 500 - Internal Server Error
// func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
// 	app.logError(r, err)
// 	app.errorResponse(w, http.StatusInternalServerError)
// }
//
// // 400 - Bad Request
// func (app *application) badRequestResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusBadRequest)
// }
//
// // 404 - Not Found
// func (app *application) notFoundResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusNotFound)
// }
//
// // 409 - Edit Conflict
// func (app *application) editConflictResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusConflict)
// }
//
// // 401 - Unauthorized
// func (app *application) invalidSessionResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusUnauthorized)
// }
//
// // 401 - Unauthorized
// func (app *application) authenticationRequiredResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusUnauthorized)
// }
//
// // 403 - Forbidden
// func (app *application) inactiveAccountResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusForbidden)
// }
//
// // 403 - Forbidden
// func (app *application) notPermittedResponse(w http.ResponseWriter) {
// 	app.errorResponse(w, http.StatusForbidden)
// }
