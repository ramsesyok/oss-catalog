package model

import "time"

// AuditLog は監査ログの 1 レコードを表す。
type AuditLog struct {
	ID         string
	EntityType string
	EntityID   string
	Action     string
	UserName   string
	Summary    *string
	CreatedAt  time.Time
}
