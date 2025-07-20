package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

func TestTagRepository_List(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{DB: db}

	query := regexp.QuoteMeta(`SELECT id, name, created_at FROM tags ORDER BY created_at DESC`)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow(uuid.NewString(), "db", now)
	mock.ExpectQuery(query).WillReturnRows(rows)

	tags, err := repo.List(context.Background())
	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{DB: db}

	t1 := &model.Tag{ID: uuid.NewString(), Name: "db", CreatedAt: func() *time.Time { v := time.Now(); return &v }()}

	query := regexp.QuoteMeta(`INSERT INTO tags (id, name, created_at) VALUES (?, ?, ?)`)
	mock.ExpectExec(query).WithArgs(t1.ID, t1.Name, t1.CreatedAt).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), t1)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTagRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &TagRepository{DB: db}

	id := uuid.NewString()
	query := regexp.QuoteMeta(`DELETE FROM tags WHERE id = ?`)
	mock.ExpectExec(query).WithArgs(id).WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
