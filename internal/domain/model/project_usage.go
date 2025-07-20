package model

import "time"

// ProjectUsage represents OSS usage within a project.
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
