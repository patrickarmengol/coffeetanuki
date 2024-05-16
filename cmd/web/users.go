package main

import (
	"net/http"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/errs"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
)

// user signup page
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.UserCreate = &model.UserCreateInput{}
	app.render(w, r, http.StatusOK, "signup.gohtml", "base", td)
}

// user signup hx
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse and decode form
	input := &model.UserCreateInput{}
	err := app.decodePostForm(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid post form format"))
		return
	}
	td.UserCreate = input

	// try to insert
	user, err := app.services.Users.Signup(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			app.logError(r, err)
			app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", "form", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.User = user

	// TODO: create registration token for user

	// redirect to login
	w.Header().Add("HX-Redirect", "/user/login")
	w.Write([]byte("user successfully registered; redirecting to login"))
}

// user login page
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.UserLogin = &model.UserLoginInput{}
	app.render(w, r, http.StatusOK, "login.gohtml", "base", td)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse and decode form
	input := &model.UserLoginInput{}
	err := app.decodePostForm(r, input)
	if err != nil {
		app.errorResponse(w, r, errs.Errorf(errs.ERRBAD, "invalid post form format"))
		return
	}
	td.UserLogin = input

	user, err := app.services.Users.Login(r.Context(), input)
	if err != nil {
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", "form", td)
		} else {
			app.errorResponse(w, r, err)
		}
		return
	}
	td.User = user

	// TODO: should i move the session stuff to service?

	// change session token to avoid session fixation
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}

	// add user id to session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", user.ID)

	// redirect to home
	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("successfully logged in; redirecting to home"))
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// change the session token
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.errorResponse(w, r, err)
		return
	}

	// remove user id from session
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// redirect to home
	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("successfully logged out; redirecting to home"))
}

func (app *application) userAccountView(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	if user.IsAnonymous() {
		app.errorResponse(w, r, errs.Errorf(errs.ERRNOTAUTHORIZED, "not authorized"))
		return
	}

	td := app.newTemplateData(r)
	td.User = user

	app.render(w, r, http.StatusOK, "account.gohtml", "base", td)
}
