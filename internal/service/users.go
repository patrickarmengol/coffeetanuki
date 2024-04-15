package service

import (
	"context"
	"database/sql"

	"github.com/patrickarmengol/coffeetanuki/internal/dba"
	"github.com/patrickarmengol/coffeetanuki/internal/errs"
	"github.com/patrickarmengol/coffeetanuki/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		db: db,
	}
}

func (serv *UserService) Exists(ctx context.Context, id int64) (bool, error) {
	return dba.UserExists(ctx, serv.db, id)
}

func (serv *UserService) Get(ctx context.Context, id int64) (*model.UserResponse, error) {
	user, err := dba.GetUser(ctx, serv.db, id)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (serv *UserService) Signup(ctx context.Context, i *model.UserCreateInput) (*model.UserResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed: %q, %q", i.FieldErrors, i.NonFieldErrors)
	}

	ucp, err := i.ToParams()
	if err != nil {
		// conversion can fail at hashing of password, but shouldn't bc validation
		return nil, err
	}

	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := dba.CreateUser(ctx, tx, ucp)
	if err != nil {
		return nil, err
	}

	// add initial permissions
	err = dba.AddPermissionsForUser(ctx, tx, user.ID, model.PermissionCodes{"beans:read", "roasters:read"})

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	ur := user.ToResponse()

	return ur, nil
}

func (serv *UserService) Login(ctx context.Context, i *model.UserLoginInput) (*model.UserResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed")
	}

	// interact with db

	// get existing user
	udb, err := dba.GetUserByEmail(ctx, serv.db, i.Email)
	if err != nil {
		// TODO: add to nonfield errors
		return nil, err
	}

	// check plaintext against hash
	err = bcrypt.CompareHashAndPassword(udb.PasswordHash, []byte(i.PasswordPlaintext))
	if err != nil {
		// TODO: add to nonfield errors
		return nil, err
	}

	// convert to response

	ur := udb.ToResponse()

	return ur, nil // is returning the user useful?
}

func (serv *UserService) GetPermissions(ctx context.Context, id int64) (model.PermissionCodes, error) {
	pcs, err := dba.GetPermissionsForUser(ctx, serv.db, id)
	if err != nil {
		return nil, err
	}
	return pcs, nil
}
