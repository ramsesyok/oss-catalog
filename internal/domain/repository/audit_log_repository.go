package repository

import (
	"context"
	"time"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// AuditLogFilter は監査ログ検索の条件を表す。
type AuditLogFilter struct {
	EntityType *string
	EntityID   *string
	From       *time.Time
	To         *time.Time
}

// AuditLogRepository は監査ログの永続化処理を定義する。
type AuditLogRepository interface {
	Search(ctx context.Context, f AuditLogFilter) ([]model.AuditLog, error)
	Create(ctx context.Context, l *model.AuditLog) error
}
