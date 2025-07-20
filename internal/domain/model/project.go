package model

import "time"

// Project represents a delivery unit project.
type Project struct {
	ID            string
	ProjectCode   string
	Name          string
	Department    *string
	Manager       *string
	DeliveryDate  *time.Time
	Description   *string
	OssUsageCount int
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
