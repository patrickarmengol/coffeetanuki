package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/patrickarmengol/somethingsomethingcoffee/internal/dba"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/errs"
	"github.com/patrickarmengol/somethingsomethingcoffee/internal/model"
)

type BeanService struct {
	db *sql.DB
}

func NewBeanService(db *sql.DB) *BeanService {
	return &BeanService{
		db: db,
	}
}

func (serv *BeanService) Create(ctx context.Context, i *model.BeanCreateInput) (*model.BeanResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed for bean create")
	}

	bcp := i.ToParams()

	// interact with db

	bdb, err := dba.CreateBean(ctx, serv.db, bcp)
	if err != nil {
		// TODO: think about how this can be improved
		if errs.ErrorCode(err) == errs.ERRUNPROCESSABLE {
			i.FieldErrors["roaster_id"] = "this ID doesn't exist or is invalid"
			return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed")
		}
		return nil, fmt.Errorf("bean repository - create: %w", err)
	}

	// convert to response

	br := bdb.ToResponse()

	return br, nil
}

func (serv *BeanService) Get(ctx context.Context, id int64) (*model.BeanResponse, error) {
	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bdb, err := dba.GetBean(ctx, tx, id)
	if err != nil {
		return nil, fmt.Errorf("bean dba - get: %w", err)
	}

	err = dba.AttachBeanAssociations(ctx, tx, bdb)
	if err != nil {
		return nil, fmt.Errorf("bean dba - get: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	br := bdb.ToResponse()

	return br, nil
}

func (serv *BeanService) Find(ctx context.Context, i *model.BeanFilterInput) ([]*model.BeanResponse, error) {
	// validate
	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRUNPROCESSABLE, "input validation failed for bean find: %q, %q", i.FieldErrors, i.NonFieldErrors)
	}

	bfp := i.ToParams()

	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bdbs, err := dba.FindBeans(ctx, tx, bfp)
	if err != nil {
		return nil, fmt.Errorf("bean dba - find: %w", err)
	}

	err = dba.AttachManyBeanAssociations(ctx, tx, bdbs)
	if err != nil {
		return nil, fmt.Errorf("bean dba - find: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	brs := []*model.BeanResponse{}
	for _, bdb := range bdbs {
		brs = append(brs, bdb.ToResponse())
	}

	return brs, nil
}

func (serv *BeanService) Update(ctx context.Context, i *model.BeanEditInput) (*model.BeanResponse, error) {
	// validate

	i.Validate()

	if !i.Valid() {
		return nil, errs.Errorf(errs.ERRBAD, "input validation failed for bean update")
	}

	bep := i.ToParams()

	// interact with db

	tx, err := serv.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bdb, err := dba.UpdateBean(ctx, tx, bep)
	if err != nil {
		return nil, fmt.Errorf("bean repository - update: %w", err)
	}

	err = dba.AttachBeanAssociations(ctx, tx, bdb)
	if err != nil {
		return nil, fmt.Errorf("bean repository - update: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// convert to response

	br := bdb.ToResponse()

	return br, nil
}

func (serv *BeanService) Delete(ctx context.Context, id int64) error {
	// interact with db

	err := dba.DeleteBean(ctx, serv.db, id)
	if err != nil {
		return fmt.Errorf("bean repository - delete: %w", err)
	}

	return nil
}
