package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
	"github.com/ramsesyok/oss-catalog/pkg/auth"
	problem "github.com/ramsesyok/oss-catalog/pkg/response"
)

// Login issues JWT token.
func (h *Handler) Login(ctx echo.Context) error {
	var req gen.LoginJSONBody
	if err := ctx.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid body")
	}
	users, _, err := h.UserRepo.Search(ctx.Request().Context(), domrepo.UserFilter{Username: req.Username, Page: 1, Size: 1})
	if err != nil {
		return err
	}
	if len(users) == 0 || users[0].PasswordHash != req.Password {
		return problem.Unauthorized(ctx, "INVALID_CREDENTIAL", "invalid username or password")
	}
	token, exp, err := auth.GenerateToken(&users[0])
	if err != nil {
		return err
	}
	res := gen.LoginResponse{AccessToken: token, ExpiresIn: exp}
	return ctx.JSON(http.StatusOK, res)
}

// Logout is a no-op placeholder.
func (h *Handler) Logout(ctx echo.Context) error {
	return ctx.NoContent(http.StatusNoContent)
}

// GetCurrentUser returns current user info based on JWT claims.
func (h *Handler) GetCurrentUser(ctx echo.Context) error {
	claims := auth.GetClaims(ctx)
	if claims == nil {
		return problem.Unauthorized(ctx, "UNAUTHORIZED", "no claims")
	}
	u, err := h.UserRepo.Get(ctx.Request().Context(), claims.Sub)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, "user not found")
		}
		return err
	}
	res := toUser(*u)
	return ctx.JSON(http.StatusOK, res)
}
