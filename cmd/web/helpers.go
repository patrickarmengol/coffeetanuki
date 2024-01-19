package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alexedwards/flow"
	"github.com/go-playground/form/v4"
	"github.com/patrickarmengol/coffeetanuki/internal/data"
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
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		// all other errors
		return err
	}

	return nil
}

func (app *application) decodeURLQuery(r *http.Request, dst any) error {
	qs := r.URL.Query()

	err := app.formDecoder.Decode(dst, qs)
	if err != nil {
		// if dst is invalid / nil pointer, panic
		var invalidDecoderError *form.InvalidDecoderError
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

func (app *application) parseSearchQuery(r *http.Request) (*data.SearchQuery, error) {
	var input struct {
		Term string `form:"term"`
		Sort string `form:"sort"`
	}

	err := app.decodeURLQuery(r, &input)
	if err != nil {
		return nil, err
	}

	if input.Sort == "" {
		input.Sort = "id_asc"
	}
	sq := data.SearchQuery{
		Term:            input.Term,
		Sort:            input.Sort,
		SortableColumns: []string{"id", "name", "location"},
	}
	return &sq, nil
}
