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
	Beans       BeanRepository
	Permissions PermissionRepository
	Roasters    RoasterRepository
	Users       UserRepository
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Beans:       BeanRepository{DB: db},
		Permissions: PermissionRepository{DB: db},
		Roasters:    RoasterRepository{DB: db},
		Users:       UserRepository{DB: db},
	}
}
