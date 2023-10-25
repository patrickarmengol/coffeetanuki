package main

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type dbConfig struct {
	dsn          string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  time.Duration
}

// helper for opening db connection pool
func openDB(dbc dbConfig) (*sql.DB, error) {
	// create empty connection pool
	db, err := sql.Open("postgres", dbc.dsn)
	if err != nil {
		return nil, err
	}

	// configure conenction pool settings
	db.SetMaxOpenConns(dbc.maxOpenConns)
	db.SetMaxIdleConns(dbc.maxIdleConns)
	db.SetConnMaxIdleTime(dbc.maxIdleTime)

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
