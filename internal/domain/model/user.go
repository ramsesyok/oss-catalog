package model

import "github.com/ramsesyok/oss-catalog/pkg/dbtime"

// User はシステム利用ユーザー情報を表すモデル。
type User struct {
	ID           string
	Username     string
	DisplayName  *string
	Email        *string
	PasswordHash string
	Roles        []string
	Active       bool
	CreatedAt    dbtime.DBTime
	UpdatedAt    dbtime.DBTime
}
