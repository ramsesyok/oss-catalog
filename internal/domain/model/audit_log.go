package model

import "time"

// AuditLog represents a single audit event.
type AuditLog struct {
	ID         string
	EntityType string
	EntityID   string
	Action     string
	UserName   string
	Summary    *string
	CreatedAt  time.Time
}
