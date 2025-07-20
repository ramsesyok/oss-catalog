package repository

import (
	"context"
	"time"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// AuditLogFilter filters audit log search.
type AuditLogFilter struct {
	EntityType *string
	EntityID   *string
	From       *time.Time
	To         *time.Time
}

// AuditLogRepository defines DB operations for AuditLog.
type AuditLogRepository interface {
	Search(ctx context.Context, f AuditLogFilter) ([]model.AuditLog, error)
	Create(ctx context.Context, l *model.AuditLog) error
}
