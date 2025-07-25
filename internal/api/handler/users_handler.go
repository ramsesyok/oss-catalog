package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
)

// ---- Users ----

// ListUsers ユーザー一覧 (GET /users)
func (*Handler) ListUsers(ctx echo.Context, params gen.ListUsersParams) error {
	return ctx.JSON(http.StatusOK, map[string]any{"placeholder": "todo"})
}

// CreateUser ユーザー作成 (POST /users)
func (*Handler) CreateUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusCreated, map[string]any{"placeholder": "todo"})
}

// GetUser ユーザー詳細 (GET /users/{userId})
func (*Handler) GetUser(ctx echo.Context, userId openapi_types.UUID) error {
	return ctx.JSON(http.StatusOK, map[string]any{"placeholder": "todo"})
}

// UpdateUser ユーザー更新 (PATCH /users/{userId})
func (*Handler) UpdateUser(ctx echo.Context, userId openapi_types.UUID) error {
	return ctx.JSON(http.StatusOK, map[string]any{"placeholder": "todo"})
}

// DeleteUser ユーザー削除 (DELETE /users/{userId})
func (*Handler) DeleteUser(ctx echo.Context, userId openapi_types.UUID) error {
	return ctx.NoContent(http.StatusNoContent)
}

// GetCurrentUser 現在ログイン中ユーザー取得 (GET /me)
func (*Handler) GetCurrentUser(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]any{"placeholder": "todo"})
}
