package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestOssComponentTagRepository_ListByOssID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentTagRepository{DB: db}

	ossID := uuid.NewString()
	query := regexp.QuoteMeta(`SELECT tg.id, tg.name, tg.created_at FROM tags tg JOIN oss_component_tags ct ON ct.tag_id = tg.id WHERE ct.oss_id = ? ORDER BY tg.created_at DESC`)
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(uuid.NewString(), "db", now)
	mock.ExpectQuery(query).WithArgs(ossID).WillReturnRows(rows)

	tags, err := repo.ListByOssID(context.Background(), ossID)
	require.NoError(t, err)
	require.Len(t, tags, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestOssComponentTagRepository_Replace(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &OssComponentTagRepository{DB: db}

	ossID := uuid.NewString()
	tagIDs := []string{uuid.NewString(), uuid.NewString()}

	mock.ExpectBegin()
	delQuery := regexp.QuoteMeta(`DELETE FROM oss_component_tags WHERE oss_id = ?`)
	mock.ExpectExec(delQuery).WithArgs(ossID).WillReturnResult(sqlmock.NewResult(1, 1))
	insQuery := regexp.QuoteMeta(`INSERT INTO oss_component_tags (oss_id, tag_id) VALUES (?, ?)`)
	for _, id := range tagIDs {
		mock.ExpectExec(insQuery).WithArgs(ossID, id).WillReturnResult(sqlmock.NewResult(1, 1))
	}
	mock.ExpectCommit()

	err = repo.Replace(context.Background(), ossID, tagIDs)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
