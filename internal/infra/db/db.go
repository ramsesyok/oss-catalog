package db

import (
	"context"
	"database/sql"
)

// DB wraps sql.DB and provides transaction management.
type DB struct {
	*sql.DB
}

// WithinTx executes fn within a transaction.
func (d *DB) WithinTx(ctx context.Context, fn func(ctx context.Context, tx *sql.Tx) error) error {
	tx, err := d.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(ctx, tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}
