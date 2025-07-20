package handler

// scope_policy_handler.go - /scope に関するハンドラ処理

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	"github.com/ramsesyok/oss-catalog/internal/domain/model"
)

// 現行スコープポリシー取得
// (GET /scope/policy)
func (h *Handler) GetScopePolicy(ctx echo.Context) error {
	p, err := h.ScopePolicyRepo.Get(ctx.Request().Context())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "policy not found")
		}
		return err
	}
	uid := uuid.MustParse(p.ID)
	res := gen.ScopePolicy{
		Id:                            &uid,
		RuntimeRequiredDefaultInScope: &p.RuntimeRequiredDefaultInScope,
		ServerEnvIncluded:             &p.ServerEnvIncluded,
		AutoMarkForksInScope:          &p.AutoMarkForksInScope,
		UpdatedAt:                     &p.UpdatedAt,
		UpdatedBy:                     &p.UpdatedBy,
	}
	return ctx.JSON(http.StatusOK, res)
}

// スコープポリシー更新 (管理者)
// (PATCH /scope/policy)
func (h *Handler) UpdateScopePolicy(ctx echo.Context) error {
	var req gen.ScopePolicyUpdateRequest
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	now := time.Now()
	id := uuid.NewString()
	p := &model.ScopePolicy{
		ID:                            id,
		RuntimeRequiredDefaultInScope: req.RuntimeRequiredDefaultInScope != nil && *req.RuntimeRequiredDefaultInScope,
		ServerEnvIncluded:             req.ServerEnvIncluded != nil && *req.ServerEnvIncluded,
		AutoMarkForksInScope:          req.AutoMarkForksInScope != nil && *req.AutoMarkForksInScope,
		UpdatedAt:                     now,
		UpdatedBy:                     "api-user",
	}
	if err := h.ScopePolicyRepo.Update(ctx.Request().Context(), p); err != nil {
		return err
	}
	uid2 := uuid.MustParse(p.ID)
	res := gen.ScopePolicy{
		Id:                            &uid2,
		RuntimeRequiredDefaultInScope: &p.RuntimeRequiredDefaultInScope,
		ServerEnvIncluded:             &p.ServerEnvIncluded,
		AutoMarkForksInScope:          &p.AutoMarkForksInScope,
		UpdatedAt:                     &p.UpdatedAt,
		UpdatedBy:                     &p.UpdatedBy,
	}
	return ctx.JSON(http.StatusOK, res)
}
