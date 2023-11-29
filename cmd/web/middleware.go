package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve existing user id from session if it exists
		id := app.sessionManager.GetInt64(r.Context(), "authenticatedUserID")
		if id == 0 {
			// user not authenticated; use anon and go to next handler
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// check if user with id exists in db
		user, err := app.repositories.Users.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				// TODO: should i force a reset on the token?
				app.invalidSessionResponse(w)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// add user to request context for use in other middleware
		r = app.contextSetUser(r, user)

		// go to next handler
		next.ServeHTTP(w, r)
	})
}

// TODO: how does this handle htmx requests?
// should i redirect to login on failed authentication required?
// should i send a HX-Redirect header?

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w)
			return
		}

		// don't store auth-required pages to browser cache
		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireActivatedUser(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if !user.Activated {
			app.inactiveAccountResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		permissions, err := app.repositories.Permissions.GetAllForUser(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		if !permissions.Contains(code) {
			app.notPermittedResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireActivatedUser(fn)
}
