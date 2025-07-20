package repository

import (
	"context"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// OssComponentFilter filters search results.
type OssComponentFilter struct {
	Name        string   // partial match on normalized_name
	Layers      []string // OR condition
	Tag         string   // exact match tag name
	InScopeOnly bool
	Page        int
	Size        int
}

// OssComponentRepository defines DB operations for OssComponent.
type OssComponentRepository interface {
	Search(ctx context.Context, f OssComponentFilter) ([]model.OssComponent, int, error)
	Create(ctx context.Context, c *model.OssComponent) error
}
