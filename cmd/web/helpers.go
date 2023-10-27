package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
)

func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
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
