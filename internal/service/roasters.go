package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/patrickarmengol/coffeetanuki/internal/dba"
	"github.com/patrickarmengol/coffeetanuki/internal/errs"
	"github.com/patrickarmengol/coffeetanuki/internal/model"
)

type RoasterService struct {
	db *sql.DB
}

func NewRoasterService(db *sql.DB) *RoasterService {
	return &RoasterService{
		db: db,
	}
}

func (serv *RoasterService) Create(ctx context.Context, i *model.RoasterCreateInput) (*model.RoasterResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed")
	}

	rcp := i.ToParams()

	// interact with db

	rdb, err := dba.CreateRoaster(ctx, serv.db, rcp)
	if err != nil {
		return nil, fmt.Errorf("roaster dba - create: %w", err)
	}

	// convert to response

	rr := rdb.ToResponse()

	return rr, nil
}

func (serv *RoasterService) Get(ctx context.Context, id int64) (*model.RoasterResponse, error) {
	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rdb, err := dba.GetRoaster(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("roaster dba - get: %w", err)
	}

	err = dba.AttachRoasterAssociations(ctx, tx, rdb)
	if err != nil {
		return nil, fmt.Errorf("roaster dba - get: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	rr := rdb.ToResponse()

	return rr, nil
}

func (serv *RoasterService) Find(ctx context.Context, i *model.RoasterFilterInput) ([]*model.RoasterResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "invalid validation failed: %q, %q", i.FieldErrors, i.NonFieldErrors)
	}

	rfp := i.ToParams()

	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rdbs, err := dba.FindRoasters(ctx, tx, rfp)
	if err != nil {
		return nil, fmt.Errorf("roaster dba - find: %w", err)
	}

	err = dba.AttachManyRoasterAssociations(ctx, tx, rdbs)
	if err != nil {
		return nil, fmt.Errorf("roaster dba - find: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	rrs := []*model.RoasterResponse{}
	for _, rdb := range rdbs {
		rrs = append(rrs, rdb.ToResponse())
	}

	return rrs, nil
}

func (serv *RoasterService) Update(ctx context.Context, i *model.RoasterEditInput) (*model.RoasterResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed")
	}

	rep := i.ToParams()

	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rdb, err := dba.UpdateRoaster(ctx, tx, rep)
	if err != nil {
		return nil, fmt.Errorf("roaster repository - update: %w", err)
	}

	err = dba.AttachRoasterAssociations(ctx, tx, rdb)
	if err != nil {
		return nil, fmt.Errorf("roaster repository - update: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	rr := rdb.ToResponse()

	return rr, nil
}

func (serv *RoasterService) Delete(ctx context.Context, id int64) error {
	// interact with db

	err := dba.DeleteRoaster(ctx, serv.db, id)
	if err != nil {
		return fmt.Errorf("roaster repository - delete: %w", err)
	}

	return nil
}
