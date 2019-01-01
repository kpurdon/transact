[![travis ci status](https://travis-ci.com/kpurdon/transact.svg?branch=master)](https://travis-ci.com/kpurdon/transact)
[![codecov status](https://codecov.io/gh/kpurdon/transact/branch/master/graph/badge.svg)](https://codecov.io/gh/kpurdon/transact)
[![godoc](https://godoc.org/github.com/kpurdon/transact?status.svg)](http://godoc.org/github.com/kpurdon/transact)
[![Go Report Card](https://goreportcard.com/badge/github.com/kpurdon/transact)](https://goreportcard.com/report/github.com/kpurdon/transact)

transact
-----

A simple helper for executing SQL queries inside a transaction and automatically handling rollback and commit scenarios.

## Example

``` go
err := transact.DoContext(ctx, db.DB, func(tx *sql.Tx) error {
    q := `SELECT ...` // any sql query that needs to run in a transaction

    _, err := tx.Exec(ctx, q, someID, count)
    if err != nil {
        return err // any errors will roll the transaction back before returning
    }

    panic("WHAT") // any panic will be recovered, roll the transaction back, and re-panic

    return nil // nil errors will commit the transaction
})
if err != nil {
    return err
}
```
