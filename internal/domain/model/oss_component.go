package model

import "time"

// OssComponent represents OSS component entity.
type OssComponent struct {
	ID               string
	Name             string
	NormalizedName   string
	HomepageURL      *string
	RepositoryURL    *string
	Description      *string
	PrimaryLanguage  *string
	Layers           []string
	DefaultUsageRole *string
	Deprecated       bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Tags             []Tag
}

// Tag represents a classification tag.
type Tag struct {
	ID        string
	Name      string
	CreatedAt *time.Time
}
