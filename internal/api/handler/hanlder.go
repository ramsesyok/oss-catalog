package handler

// hanlder.go - ハンドラ系の共通部分
// (エンドポイント定義無し)

import (
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

type Handler struct {
	AuditRepo             domrepo.AuditLogRepository
	ScopePolicyRepo       domrepo.ScopePolicyRepository
	OssComponentRepo      domrepo.OssComponentRepository
	OssComponentLayerRepo domrepo.OssComponentLayerRepository
	OssComponentTagRepo   domrepo.OssComponentTagRepository
	OssVersionRepo        domrepo.OssVersionRepository
	ProjectRepo           domrepo.ProjectRepository
	ProjectUsageRepo      domrepo.ProjectUsageRepository
}
