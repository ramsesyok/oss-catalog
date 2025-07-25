package model

import "time"

// User はシステム利用ユーザー情報を表すモデル。
type User struct {
	ID           string
	Username     string
	DisplayName  *string
	Email        *string
	PasswordHash string
	Roles        []string
	Active       bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
