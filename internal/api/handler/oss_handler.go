package handler

// oss_handler.go - /oss に関するハンドラ処理

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
)

// OSSコンポーネント一覧取得
// (GET /oss)
func (*Handler) ListOssComponents(ctx echo.Context, params gen.ListOssComponentsParams) error {
	return nil
}

// OSSコンポーネント作成
// (POST /oss)
func (*Handler) CreateOssComponent(ctx echo.Context) error {
	return nil
}

// OSSコンポーネントを非推奨 (deprecated=true) に設定
// (DELETE /oss/{ossId})
func (*Handler) DeprecateOssComponent(ctx echo.Context, ossId openapi_types.UUID) error {
	return nil
}

// OSSコンポーネント詳細
// (GET /oss/{ossId})
func (*Handler) GetOssComponent(ctx echo.Context, ossId openapi_types.UUID) error {
	return nil
}

// OSSコンポーネント更新 (部分)
// (PATCH /oss/{ossId})
func (*Handler) UpdateOssComponent(ctx echo.Context, ossId openapi_types.UUID) error {
	return nil
}

// 指定 OSS のバージョン一覧
// (GET /oss/{ossId}/versions)
func (*Handler) ListOssVersions(ctx echo.Context, ossId openapi_types.UUID, params gen.ListOssVersionsParams) error {
	return nil
}

// バージョン追加
// (POST /oss/{ossId}/versions)
func (*Handler) CreateOssVersion(ctx echo.Context, ossId openapi_types.UUID) error {
	return nil
}

// バージョン削除 (論理/物理は実装方針による)
// (DELETE /oss/{ossId}/versions/{versionId})
func (*Handler) DeleteOssVersion(ctx echo.Context, ossId openapi_types.UUID, versionId openapi_types.UUID) error {
	return nil
}

// バージョン詳細
// (GET /oss/{ossId}/versions/{versionId})
func (*Handler) GetOssVersion(ctx echo.Context, ossId openapi_types.UUID, versionId openapi_types.UUID) error {
	return nil
}

// バージョン更新
// (PATCH /oss/{ossId}/versions/{versionId})
func (*Handler) UpdateOssVersion(ctx echo.Context, ossId openapi_types.UUID, versionId openapi_types.UUID) error {
	return nil
}
