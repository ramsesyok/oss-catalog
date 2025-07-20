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
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

func TestAuditLogRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &AuditLogRepository{DB: db}

	f := domrepo.AuditLogFilter{EntityType: func() *string { s := "PROJECT"; return &s }()}

	query := regexp.QuoteMeta("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs WHERE entity_type = ? ORDER BY created_at DESC")
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "entity_type", "entity_id", "action", "user_name", "summary", "created_at"}).
		AddRow(uuid.NewString(), "PROJECT", "1", "CREATE", "user", "created", now)
	mock.ExpectQuery(query).WithArgs("PROJECT").WillReturnRows(rows)

	logs, err := repo.Search(context.Background(), f)
	require.NoError(t, err)
	require.Len(t, logs, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestAuditLogRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := &AuditLogRepository{DB: db}

	l := &model.AuditLog{
		ID:         uuid.NewString(),
		EntityType: "PROJECT",
		EntityID:   "1",
		Action:     "CREATE",
		UserName:   "user",
		CreatedAt:  time.Now(),
	}

	query := regexp.QuoteMeta("INSERT INTO audit_logs (id, entity_type, entity_id, action, user_name, summary, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(query).
		WithArgs(l.ID, l.EntityType, l.EntityID, l.Action, l.UserName, l.Summary, l.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(context.Background(), l)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
