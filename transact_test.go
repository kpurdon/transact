package transact

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	sqlmock "github.com/data-dog/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDoContext(t *testing.T) {
	testCases := []struct {
		label  string
		setup  func(t *testing.T) (context.Context, context.CancelFunc)
		txFunc func(tx *sql.Tx) error
		check  func(t *testing.T, ctx context.Context, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error)
	}{
		{
			label: "an error is returned on any error",
			setup: func(t *testing.T) (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			txFunc: func(tx *sql.Tx) error {
				return assert.AnError
			},
			check: func(t *testing.T, ctx context.Context, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
				mock.ExpectBegin()
				mock.ExpectRollback()
				require.EqualError(t, DoContext(ctx, db, txFunc), assert.AnError.Error())
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		}, {
			label: "a panic is raised on any panic",
			setup: func(t *testing.T) (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			txFunc: func(tx *sql.Tx) error {
				panic(assert.AnError)
			},
			check: func(t *testing.T, ctx context.Context, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
				mock.ExpectBegin()
				mock.ExpectRollback()
				require.Panics(t, func() { DoContext(ctx, db, txFunc) })
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		}, {
			label: "no error is returned on any success",
			setup: func(t *testing.T) (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			txFunc: func(tx *sql.Tx) error {
				return nil
			},
			check: func(t *testing.T, ctx context.Context, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
				mock.ExpectBegin()
				mock.ExpectCommit()
				require.NoError(t, DoContext(ctx, db, txFunc))
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
		// TODO: figure out context cancellation ... maybe just the sqlmock driver
		// {
		// 	label: "a context cancellation cancels the transaction",
		// 	setup: func(t *testing.T) (context.Context, context.CancelFunc) {
		// 		ctx := context.Background()
		// 		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		// 		return ctx, cancel
		// 	},
		// 	txFunc: func(tx *sql.Tx) error {
		// 		time.Sleep(10 * time.Second) // anything longer than the timeout
		// 		return nil
		// 	},
		// 	check: func(t *testing.T, ctx context.Context, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
		// 		mock.ExpectBegin()
		// 		mock.ExpectRollback()
		// 		require.Error(t, DoContext(ctx, db, txFunc))
		// 		assert.NoError(t, mock.ExpectationsWereMet())
		// 	},
		// },
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.label)

			ctx, cancel := tc.setup(t)
			defer cancel()

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.check(t, ctx, db, mock, tc.txFunc)
		})
	}
}

func TestDo(t *testing.T) {
	testCases := []struct {
		label  string
		txFunc func(tx *sql.Tx) error
		check  func(t *testing.T, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error)
	}{
		{
			label: "an error is returned on any error",
			txFunc: func(tx *sql.Tx) error {
				return assert.AnError
			},
			check: func(t *testing.T, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
				mock.ExpectBegin()
				mock.ExpectRollback()
				require.EqualError(t, Do(db, txFunc), assert.AnError.Error())
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		}, {
			label: "a panic is raised on any panic",
			txFunc: func(tx *sql.Tx) error {
				panic(assert.AnError)
			},
			check: func(t *testing.T, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
				mock.ExpectBegin()
				mock.ExpectRollback()
				require.Panics(t, func() { Do(db, txFunc) })
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		}, {
			label: "no error is returned on any success",
			txFunc: func(tx *sql.Tx) error {
				return nil
			},
			check: func(t *testing.T, db *sql.DB, mock sqlmock.Sqlmock, txFunc func(tx *sql.Tx) error) {
				mock.ExpectBegin()
				mock.ExpectCommit()
				require.NoError(t, Do(db, txFunc))
				assert.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for i, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			t.Log(tc.label)

			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tc.check(t, db, mock, tc.txFunc)
		})
	}
}
