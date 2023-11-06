package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound   = errors.New("record not found")
	ErrEditConflict     = errors.New("edit conflict")
	ErrInvalidRoasterID = errors.New("invalid roaster id")
)

type Repositories struct {
	Roasters RoasterRepository
	Beans    BeanRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Roasters: RoasterRepository{DB: db},
		Beans:    BeanRepository{DB: db},
	}
}
