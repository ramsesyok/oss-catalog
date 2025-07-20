package handler

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
)

type Handler struct {
}

// 監査ログ簡易検索 (Phase1簡易)
// (GET /audit)
func (*Handler) SearchAuditLogs(ctx echo.Context, params gen.SearchAuditLogsParams) error {
	return nil
}

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

// プロジェクト一覧
// (GET /projects)
func (*Handler) ListProjects(ctx echo.Context, params gen.ListProjectsParams) error {
	return nil
}

// プロジェクト作成
// (POST /projects)
func (*Handler) CreateProject(ctx echo.Context) error {
	return nil
}

// プロジェクト削除 (論理予定)
// (DELETE /projects/{projectId})
func (*Handler) DeleteProject(ctx echo.Context, projectId openapi_types.UUID) error {
	return nil
}

// プロジェクト詳細
// (GET /projects/{projectId})
func (*Handler) GetProject(ctx echo.Context, projectId openapi_types.UUID) error {
	return nil
}

// プロジェクト更新
// (PATCH /projects/{projectId})
func (*Handler) UpdateProject(ctx echo.Context, projectId openapi_types.UUID) error {
	return nil
}

// プロジェクト納品用エクスポート (プレースホルダ)
// (GET /projects/{projectId}/export)
func (*Handler) ExportProjectArtifacts(ctx echo.Context, projectId openapi_types.UUID, params gen.ExportProjectArtifactsParams) error {
	return nil
}

// プロジェクト中利用 OSS 一覧
// (GET /projects/{projectId}/usages)
func (*Handler) ListProjectUsages(ctx echo.Context, projectId openapi_types.UUID, params gen.ListProjectUsagesParams) error {
	return nil
}

// プロジェクト利用追加
// (POST /projects/{projectId}/usages)
func (*Handler) CreateProjectUsage(ctx echo.Context, projectId openapi_types.UUID) error {
	return nil
}

// 利用削除
// (DELETE /projects/{projectId}/usages/{usageId})
func (*Handler) DeleteProjectUsage(ctx echo.Context, projectId openapi_types.UUID, usageId openapi_types.UUID) error {
	return nil
}

// 利用情報更新
// (PATCH /projects/{projectId}/usages/{usageId})
func (*Handler) UpdateProjectUsage(ctx echo.Context, projectId openapi_types.UUID, usageId openapi_types.UUID) error {
	return nil
}

// スコープ判定更新
// (PATCH /projects/{projectId}/usages/{usageId}/scope)
func (*Handler) UpdateProjectUsageScope(ctx echo.Context, projectId openapi_types.UUID, usageId openapi_types.UUID) error {
	return nil
}

// 現行スコープポリシー取得
// (GET /scope/policy)
func (*Handler) GetScopePolicy(ctx echo.Context) error {
	return nil
}

// スコープポリシー更新 (管理者)
// (PATCH /scope/policy)
func (*Handler) UpdateScopePolicy(ctx echo.Context) error {
	return nil
}

// タグ一覧
// (GET /tags)
func (*Handler) ListTags(ctx echo.Context) error {
	return nil
}

// タグ作成
// (POST /tags)
func (*Handler) CreateTag(ctx echo.Context) error {
	return nil
}

// タグ削除
// (DELETE /tags/{tagId})
func (*Handler) DeleteTag(ctx echo.Context, tagId openapi_types.UUID) error {
	return nil
}
