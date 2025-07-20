package repository

import (
	"context"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// ScopePolicyRepository defines DB operations for ScopePolicy.
type ScopePolicyRepository interface {
	// Get returns current scope policy. Returns nil if not found.
	Get(ctx context.Context) (*model.ScopePolicy, error)
	// Update upserts the scope policy record.
	Update(ctx context.Context, p *model.ScopePolicy) error
}
