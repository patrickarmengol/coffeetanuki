package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Roasters RoasterModel
	Beans    BeanModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Roasters: RoasterModel{DB: db},
		Beans:    BeanModel{DB: db},
	}
}
