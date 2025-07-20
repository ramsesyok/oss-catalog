package repository

import (
	"context"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// ProjectFilter filters project listing.
type ProjectFilter struct {
	Code string
	Name string
	Page int
	Size int
}

// ProjectRepository defines DB operations for Project.
type ProjectRepository interface {
	Search(ctx context.Context, f ProjectFilter) ([]model.Project, int, error)
	Get(ctx context.Context, id string) (*model.Project, error)
	Create(ctx context.Context, p *model.Project) error
	Update(ctx context.Context, p *model.Project) error
	Delete(ctx context.Context, id string) error
}
