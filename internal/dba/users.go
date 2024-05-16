package dba

import (
	"context"
	"database/sql"
	"errors"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/errs"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
)

func CreateUser(ctx context.Context, dbtx DBTX, p *model.UserCreateParams) (*model.UserDB, error) {
	stmt := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version
	`

	args := []any{p.Name, p.Email, p.PasswordHash, p.Activated}

	user := model.UserDB{
		Name:         p.Name,
		Email:        p.Email,
		PasswordHash: p.PasswordHash,
		Activated:    p.Activated,
	}

	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, errDuplicate("users", "email", user.Email)
		default:
			return nil, err
		}
	}

	return &user, nil
}

func UserExists(ctx context.Context, dbtx DBTX, id int64) (bool, error) {
	stmt := `
	SELECT EXISTS(SELECT true FROM users WHERE id = $1)
	`

	args := []any{id}

	var exists bool
	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&exists)
	return exists, err
}

func GetUserByEmail(ctx context.Context, dbtx DBTX, email string) (*model.UserDB, error) {
	stmt := `
	SELECT id, name, email, password_hash, activated, created_at, version
	FROM users
	WHERE email = $1
	`

	args := []any{email}

	var user model.UserDB
	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Activated, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errs.Errorf(errs.ERRNOTFOUND, "record with email %s not found", email)
		default:
			return nil, err
		}
	}

	return &user, nil
}

func GetUser(ctx context.Context, dbtx DBTX, id int64) (*model.UserDB, error) {
	stmt := `
	SELECT id, name, email, password_hash, activated, created_at, version
	FROM users
	WHERE id = $1
	`

	args := []any{id}

	var user model.UserDB
	err := dbtx.QueryRowContext(ctx, stmt, args...).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Activated, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, errRecordNotFound("users", id)
		default:
			return nil, err
		}
	}

	return &user, nil
}
