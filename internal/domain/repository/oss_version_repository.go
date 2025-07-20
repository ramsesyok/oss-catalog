package repository

import (
	"context"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// OssVersionFilter filters search results.
type OssVersionFilter struct {
	OssID        string
	ReviewStatus string
	ScopeStatus  string
	Page         int
	Size         int
}

// OssVersionRepository defines DB operations for OssVersion.
type OssVersionRepository interface {
	Search(ctx context.Context, f OssVersionFilter) ([]model.OssVersion, int, error)
	Get(ctx context.Context, id string) (*model.OssVersion, error)
	Create(ctx context.Context, v *model.OssVersion) error
	Update(ctx context.Context, v *model.OssVersion) error
	Delete(ctx context.Context, id string) error
}
