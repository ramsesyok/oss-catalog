package repository

import (
	"context"
	"time"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// ProjectUsageFilter filters usage listing.
type ProjectUsageFilter struct {
	ProjectID   string
	ScopeStatus string
	UsageRole   string
	Direct      *bool
	Page        int
	Size        int
}

// ProjectUsageRepository defines DB operations for ProjectUsage.
type ProjectUsageRepository interface {
	Search(ctx context.Context, f ProjectUsageFilter) ([]model.ProjectUsage, int, error)
	Create(ctx context.Context, u *model.ProjectUsage) error
	Update(ctx context.Context, u *model.ProjectUsage) error
	Delete(ctx context.Context, id string) error
	UpdateScope(ctx context.Context, id string, scopeStatus string, inclusionNote *string, evaluatedAt time.Time, evaluatedBy *string) error
}
