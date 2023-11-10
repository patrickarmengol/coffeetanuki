package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repositories struct {
	Roasters RoasterRepository
	Beans    BeanRepository
	Users    UserRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Roasters: RoasterRepository{DB: db},
		Beans:    BeanRepository{DB: db},
		Users:    UserRepository{DB: db},
	}
}
