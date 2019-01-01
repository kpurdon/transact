// Package transact provides a simple helper for running sql queries inside a transaction without
// needing to handle all rollback and commit scenarios manually.
package transact

import (
	"context"
	"database/sql"
)

// DoContext executes the given txFunc inside of a new transaction handling all possible rollback
// and commit scenarios.
func DoContext(ctx context.Context, db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return
	}

	defer func() {
		if pErr := recover(); pErr != nil {
			tx.Rollback()
			panic(pErr)
		}

		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}

	}()

	return txFunc(tx)
}

// Do executes the given txFunc inside of a new transaction handling all possible rollback and
// commit scenarios.
func Do(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	return DoContext(context.Background(), db, txFunc)
}
