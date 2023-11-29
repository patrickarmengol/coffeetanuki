package main

import (
	"errors"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/data"
	"github.com/patrickarmengol/coffeetanuki/internal/validator"
)

type userSignupForm struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type userLoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.Validator = validator.New()
	td.User = &data.User{}
	app.render(w, r, http.StatusOK, "signup.gohtml", "base", td)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse and decode form
	var form userSignupForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass into model
	user := &data.User{
		Name:      form.Name,
		Email:     form.Email,
		Activated: false,
	}
	td.User = user
	err = user.Password.Set(form.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// validate
	v := validator.New()
	td.Validator = v
	user.Validate(v)

	// invalid form input
	if !v.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", "form", td)
		return
	}

	// try to insert
	err = app.repositories.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddFieldError("email", "a user with this email address already exists")
			app.render(w, r, http.StatusUnprocessableEntity, "signup.gohtml", "form", td)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// set permissions for user
	err = app.repositories.Permissions.AddForUser(user.ID, "beans:read", "roasters:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: create registration token for user

	// redirect to login
	w.Header().Add("HX-Redirect", "/user/login")
	w.Write([]byte("user successfully registered; redirecting to login"))
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// render form with empty model
	td.Validator = validator.New()
	td.User = &data.User{}
	app.render(w, r, http.StatusOK, "login.gohtml", "base", td)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	td := app.newTemplateData(r)

	// parse and decode form
	var form userLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w)
		return
	}

	// pass into model
	td.User = &data.User{
		Email: form.Email,
	}

	// validate email and pass individually instead of full user
	v := validator.New()
	td.Validator = v
	data.ValidateEmail(v, form.Email)
	data.ValidatePasswordPlaintext(v, form.Password)

	// invalid form input
	if !v.Valid() {
		app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", "form", td)
		return
	}

	// check for db user matching email
	dbUser, err := app.repositories.Users.GetByEmail(form.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddNonFieldError("email or password invalid")
			app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", "form", td)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// check for db password hash matching form password
	match, err := dbUser.Password.Matches(form.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		v.AddNonFieldError("email or password invalid")
		app.render(w, r, http.StatusUnprocessableEntity, "login.gohtml", "form", td)
	}

	// change session token to avoid session fixation
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// add user id to session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", dbUser.ID)

	// redirect to home
	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("successfully logged in; redirecting to home"))
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// change the session token
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// remove user id from session
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// redirect to home
	w.Header().Add("HX-Redirect", "/")
	w.Write([]byte("successfully logged out; redirecting to home"))
}
