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
		Beans:       BeanRepository{db: db},
		Permissions: PermissionRepository{db: db},
		Roasters:    RoasterRepository{db: db},
		Users:       UserRepository{db: db},
	}
}
