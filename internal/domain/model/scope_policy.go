package model

import "time"

// ScopePolicy represents automatic scope determination policy.
type ScopePolicy struct {
	ID                            string
	RuntimeRequiredDefaultInScope bool
	ServerEnvIncluded             bool
	AutoMarkForksInScope          bool
	UpdatedAt                     time.Time
	UpdatedBy                     string
}
