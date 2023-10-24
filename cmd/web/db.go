package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

// helper for opening db connection pool
func (app *application) openDB() (*sql.DB, error) {
	// create empty connection pool
	db, err := sql.Open("postgres", app.config.db.dsn)
	if err != nil {
		return nil, err
	}

	// configure conenction pool settings
	db.SetMaxOpenConns(app.config.db.maxOpenConns)
	db.SetMaxIdleConns(app.config.db.maxIdleConns)
	db.SetConnMaxIdleTime(app.config.db.maxIdleTime)

	// test connection with a ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// return pool
	return db, nil
}
