package migration

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

func TestApply_SQLite(t *testing.T) {
	db, err := sql.Open("sqlite3", "file:test?mode=memory&cache=shared")
	require.NoError(t, err)
	defer db.Close()

	err = Apply(db, "file:test?mode=memory&cache=shared")
	require.NoError(t, err)

	var name string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='tags'`).Scan(&name)
	require.NoError(t, err)
	require.Equal(t, "tags", name)
}
