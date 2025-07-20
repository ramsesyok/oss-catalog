package db

import (
	"context"
	"database/sql"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// DB wraps sql.DB and provides transaction management.
type DB struct {
	*sql.DB
}

// Open opens a database connection using the provided DSN.
// If the DSN starts with "postgres://" it uses the postgres driver,
// otherwise it falls back to SQLite.
func Open(dsn string) (*DB, error) {
	driver := "sqlite3"
	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		driver = "postgres"
	}
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return &DB{DB: db}, nil
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
