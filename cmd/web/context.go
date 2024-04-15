package main

import (
	"context"
	"net/http"

	"github.com/patrickarmengol/coffeetanuki/internal/model"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) contextSetUser(r *http.Request, user *model.UserResponse) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *model.UserResponse {
	user, ok := r.Context().Value(userContextKey).(*model.UserResponse)
	if !ok {
		panic("missing user value in request context")
	}

	return user
}
