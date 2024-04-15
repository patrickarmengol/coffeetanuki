package dba

import (
	"context"
	"database/sql"
)

type DBTX interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

// reference for tx pattern
//
// tx, err := repo.db.BeginTx(ctx, nil)
// if err != nil {
// 	return nil, err
// }
// defer tx.Rollback()
//
// user, err := getUser(ctx, tx, id)
// if err != nil {
// 	return nil, err
// }
//
// err = tx.Commit()
// if err != nil {
// 	return nil, err
// }
//
// return user, nil
