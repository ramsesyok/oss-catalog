package repository

import "context"

// OssComponentLayerRepository defines operations on oss_component_layers table.
type OssComponentLayerRepository interface {
	// ListByOssID returns layers associated with a component ordered by layer.
	ListByOssID(ctx context.Context, ossID string) ([]string, error)
	// Replace replaces layers for a component with given layers.
	Replace(ctx context.Context, ossID string, layers []string) error
}
