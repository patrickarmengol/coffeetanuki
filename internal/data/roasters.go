package data

import (
	"database/sql"
	"time"
)

type Roaster struct {
	ID        int64
	Name      string
	Location  string
	CreatedAt time.Time
	Version   int
}

type RoasterModel struct {
	DB *sql.DB
}
