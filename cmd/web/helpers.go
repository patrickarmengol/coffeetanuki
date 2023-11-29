package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/flow"
	"github.com/go-playground/form/v4"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(flow.Param(r.Context(), "id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (app *application) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		// if dst is invalid / nil pointer, panic
		var invalidDecoderError *form.InvalidEncodeError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// all other errors
		return err
	}

	return nil
}

func (app *application) isAuthenticated(r *http.Request) bool {
	user := app.contextGetUser(r)
	return !user.IsAnonymous()
}
