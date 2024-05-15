package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/patrickarmengol/coffeetanuki/internal/errs"
	"github.com/patrickarmengol/coffeetanuki/internal/model"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// retrieve existing user id from session if it exists
		id := app.sessionManager.GetInt64(r.Context(), "authenticatedUserID")
		if id == 0 {
			// user not authenticated; use anon and go to next handler
			r = app.contextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// check if user with id exists in db
		user, err := app.services.Users.Get(r.Context(), id)
		if err != nil {
			switch {
			case errs.ErrorCode(err) == errs.ERRNOTFOUND:
				// TODO: should i force a reset on the token?
				app.errorResponse(w, r, errs.Errorf(errs.ERRNOTFOUND, "user associated with session token not found"))
			default:
				app.errorResponse(w, r, err)
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
			app.errorResponse(w, r, errs.Errorf(errs.ERRNOTAUTHORIZED, "user authentication required"))
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
			app.errorResponse(w, r, errs.Errorf(errs.ERRNOTAUTHORIZED, "user account inactive"))
			return
		}

		next.ServeHTTP(w, r)
	})

	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			obj, act, found := strings.Cut(code, ":")
			if !found {
				panic("permission codes should have a semicolon separator")
			}

			sub := app.contextGetUser(r).Name

			ok, err := app.rbacEnforcer.Enforce(sub, obj, act)
			if err != nil {
				app.errorResponse(w, r, err)
				return
			}
			if !ok {
				app.errorResponse(w, r, errs.Errorf(errs.ERRNOTAUTHORIZED, "user does not have required permissions"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func setTimeout(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second) // TODO: make timeout configurable
		defer cancel()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
