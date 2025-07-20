package model

import "time"

// ScopePolicy は自動スコープ判定のポリシーを表す。
type ScopePolicy struct {
	ID                            string
	RuntimeRequiredDefaultInScope bool
	ServerEnvIncluded             bool
	AutoMarkForksInScope          bool
	UpdatedAt                     time.Time
	UpdatedBy                     string
}
