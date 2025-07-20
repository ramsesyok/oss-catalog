package model

import "time"

// OssComponent は OSS コンポーネントを表すドメインモデル。
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

// Tag はコンポーネントに付与される分類タグ。
type Tag struct {
	ID        string
	Name      string
	CreatedAt *time.Time
}
