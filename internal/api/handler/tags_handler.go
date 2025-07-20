package handler

// tags_handler.go - /tags に関するハンドラ処理

import (
	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

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
