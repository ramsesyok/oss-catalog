package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

// AuditLogRepository implements domrepo.AuditLogRepository.
type AuditLogRepository struct {
	DB *sql.DB
}

var _ domrepo.AuditLogRepository = (*AuditLogRepository)(nil)

// Search returns audit logs matching filter ordered by created_at desc.
func (r *AuditLogRepository) Search(ctx context.Context, f domrepo.AuditLogFilter) ([]model.AuditLog, error) {
	var args []any
	var wheres []string
	if f.EntityType != nil && *f.EntityType != "" {
		wheres = append(wheres, "entity_type = ?")
		args = append(args, *f.EntityType)
	}
	if f.EntityID != nil && *f.EntityID != "" {
		wheres = append(wheres, "entity_id = ?")
		args = append(args, *f.EntityID)
	}
	if f.From != nil {
		wheres = append(wheres, "created_at >= ?")
		args = append(args, *f.From)
	}
	if f.To != nil {
		wheres = append(wheres, "created_at <= ?")
		args = append(args, *f.To)
	}
	whereSQL := ""
	if len(wheres) > 0 {
		whereSQL = "WHERE " + strings.Join(wheres, " AND ")
	}
	query := fmt.Sprintf("SELECT id, entity_type, entity_id, action, user_name, summary, created_at FROM audit_logs %s ORDER BY created_at DESC", whereSQL)
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		var summary sql.NullString
		if err := rows.Scan(&l.ID, &l.EntityType, &l.EntityID, &l.Action, &l.UserName, &summary, &l.CreatedAt); err != nil {
			return nil, err
		}
		if summary.Valid {
			l.Summary = &summary.String
		}
		logs = append(logs, l)
	}
	return logs, rows.Err()
}

// Create inserts a new audit log entry.
func (r *AuditLogRepository) Create(ctx context.Context, l *model.AuditLog) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO audit_logs (id, entity_type, entity_id, action, user_name, summary, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		l.ID, l.EntityType, l.EntityID, l.Action, l.UserName, l.Summary, l.CreatedAt,
	)
	return err
}
