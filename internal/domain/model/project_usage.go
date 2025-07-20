package model

import "time"

// ProjectUsage はプロジェクト内での OSS 利用状況を表す。
type ProjectUsage struct {
	ID               string
	ProjectID        string
	OssID            string
	OssVersionID     string
	UsageRole        string
	ScopeStatus      string
	InclusionNote    *string
	DirectDependency bool
	AddedAt          time.Time
	EvaluatedAt      *time.Time
	EvaluatedBy      *string
}

// ProjectUsageFilter not defined here - in repository package.
