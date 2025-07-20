package db

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDB_WithinTx_Commit(t *testing.T) {
	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer rawDB.Close()

	d := &DB{DB: rawDB}
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = d.WithinTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return nil
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDB_WithinTx_Rollback(t *testing.T) {
	rawDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer rawDB.Close()

	d := &DB{DB: rawDB}
	mock.ExpectBegin()
	mock.ExpectRollback()

	err = d.WithinTx(context.Background(), func(ctx context.Context, tx *sql.Tx) error {
		return errors.New("fail")
	})
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
