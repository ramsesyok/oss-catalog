package model

import "time"

// OssVersion は OSS コンポーネントのバージョン情報を表す。
type OssVersion struct {
	ID                      string
	OssID                   string
	Version                 string
	ReleaseDate             *time.Time
	LicenseExpressionRaw    *string
	LicenseConcluded        *string
	Purl                    *string
	CpeList                 []string
	HashSha256              *string
	Modified                bool
	ModificationDescription *string
	ReviewStatus            string
	LastReviewedAt          *time.Time
	ScopeStatus             string
	SupplierType            *string
	ForkOriginURL           *string
	CreatedAt               time.Time
	UpdatedAt               time.Time
}
