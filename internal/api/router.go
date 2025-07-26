package api

import (
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/ramsesyok/oss-catalog/internal/api/gen"
	"github.com/ramsesyok/oss-catalog/internal/api/handler"
	"github.com/ramsesyok/oss-catalog/pkg/auth"
	problem "github.com/ramsesyok/oss-catalog/pkg/response"
	"os"
)

// RegisterRoutes registers handlers with JWT auth.
func RegisterRoutes(e *echo.Echo, h *handler.Handler) {
	cfg := echojwt.Config{
		SigningKey:    []byte(os.Getenv("JWT_SECRET")),
		ContextKey:    "authUser",
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return &auth.Claims{} },
		ErrorHandler: func(c echo.Context, err error) error {
			return problem.Unauthorized(c, "UNAUTHORIZED", "token invalid or expired")
		},
	}
	authRequired := echojwt.WithConfig(cfg)

	wrapper := gen.ServerInterfaceWrapper{Handler: h}

	e.POST("/auth/login", wrapper.Login)
	e.POST("/auth/logout", wrapper.Logout, authRequired)

	g := e.Group("", authRequired)
	g.GET("/audit", wrapper.SearchAuditLogs)
	g.GET("/me", wrapper.GetCurrentUser)
	g.GET("/oss", wrapper.ListOssComponents)
	g.POST("/oss", wrapper.CreateOssComponent)
	g.DELETE("/oss/:ossId", wrapper.DeprecateOssComponent)
	g.GET("/oss/:ossId", wrapper.GetOssComponent)
	g.PATCH("/oss/:ossId", wrapper.UpdateOssComponent)
	g.GET("/oss/:ossId/versions", wrapper.ListOssVersions)
	g.POST("/oss/:ossId/versions", wrapper.CreateOssVersion)
	g.DELETE("/oss/:ossId/versions/:versionId", wrapper.DeleteOssVersion)
	g.GET("/oss/:ossId/versions/:versionId", wrapper.GetOssVersion)
	g.PATCH("/oss/:ossId/versions/:versionId", wrapper.UpdateOssVersion)
	g.GET("/projects", wrapper.ListProjects)
	g.POST("/projects", wrapper.CreateProject)
	g.DELETE("/projects/:projectId", wrapper.DeleteProject)
	g.GET("/projects/:projectId", wrapper.GetProject)
	g.PATCH("/projects/:projectId", wrapper.UpdateProject)
	g.GET("/projects/:projectId/export", wrapper.ExportProjectArtifacts)
	g.GET("/projects/:projectId/usages", wrapper.ListProjectUsages)
	g.POST("/projects/:projectId/usages", wrapper.CreateProjectUsage)
	g.DELETE("/projects/:projectId/usages/:usageId", wrapper.DeleteProjectUsage)
	g.PATCH("/projects/:projectId/usages/:usageId", wrapper.UpdateProjectUsage)
	g.PATCH("/projects/:projectId/usages/:usageId/scope", wrapper.UpdateProjectUsageScope)
	g.GET("/scope/policy", wrapper.GetScopePolicy)
	g.PATCH("/scope/policy", wrapper.UpdateScopePolicy)
	g.GET("/tags", wrapper.ListTags)
	g.POST("/tags", wrapper.CreateTag)
	g.DELETE("/tags/:tagId", wrapper.DeleteTag)
	g.GET("/users", wrapper.ListUsers)
	g.POST("/users", wrapper.CreateUser)
	g.DELETE("/users/:userId", wrapper.DeleteUser)
	g.GET("/users/:userId", wrapper.GetUser)
	g.PATCH("/users/:userId", wrapper.UpdateUser)
}
