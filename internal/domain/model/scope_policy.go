package model

import "github.com/ramsesyok/oss-catalog/pkg/dbtime"

// ScopePolicy は自動スコープ判定のポリシーを表す。
type ScopePolicy struct {
	ID                            string
	RuntimeRequiredDefaultInScope bool
	ServerEnvIncluded             bool
	AutoMarkForksInScope          bool
	UpdatedAt                     dbtime.DBTime
	UpdatedBy                     string
}
