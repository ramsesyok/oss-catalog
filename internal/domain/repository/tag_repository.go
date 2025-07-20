package repository

import (
	"context"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// TagRepository defines DB operations for Tag.
type TagRepository interface {
	// List returns all tags ordered by creation date descending.
	List(ctx context.Context) ([]model.Tag, error)
	// Create inserts a new tag record.
	Create(ctx context.Context, t *model.Tag) error
	// Delete removes a tag by ID.
	Delete(ctx context.Context, id string) error
}
