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
}

func NewModels(db *sql.DB) Models {
	return Models{
		Roasters: RoasterModel{DB: db},
	}
}
