package handler

// projects_handler.go - /projects に関するハンドラ処理

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
)

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

// プロジェクト納品用エクスポート (プレーホルダ)
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
