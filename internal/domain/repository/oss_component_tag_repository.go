package repository

import (
	"context"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// OssComponentTagRepository defines operations on oss_component_tags table.
type OssComponentTagRepository interface {
	// ListByOssID returns tags associated with a component ordered by created_at.
	ListByOssID(ctx context.Context, ossID string) ([]model.Tag, error)
	// Replace replaces tags for a component with given tagIDs.
	Replace(ctx context.Context, ossID string, tagIDs []string) error
}
