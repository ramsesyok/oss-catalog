package migration

import (
	"database/sql"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/ramsesyok/oss-catalog/migrations"
)

// Apply runs all up migrations using golang-migrate. If the schema is already up to date,
// it does nothing.
func Apply(db *sql.DB, dsn string) error {
	var (
		driver database.Driver
		err    error
	)

	if strings.HasPrefix(dsn, "postgres://") || strings.HasPrefix(dsn, "postgresql://") {
		driver, err = postgres.WithInstance(db, &postgres.Config{})
	} else {
		driver, err = sqlite3.WithInstance(db, &sqlite3.Config{})
	}
	if err != nil {
		return err
	}

	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", src, "", driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
