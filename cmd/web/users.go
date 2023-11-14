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
	td := newTemplateData()

	// render form with empty model
	td.Validator = validator.New()
	td.User = &data.User{}
	app.render(w, r, http.StatusOK, "signup.gohtml", "base", td)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	td := newTemplateData()

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

	// TODO: set permissions for user

	// TODO: create registration token for user

	// use HX-Redirect instead
	w.Header().Add("HX-Redirect", "/user/login")
	w.Write([]byte("user successfully registered; redirecting to login"))
}
