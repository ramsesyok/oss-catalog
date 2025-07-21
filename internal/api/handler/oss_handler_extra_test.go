package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	"github.com/ramsesyok/oss-catalog/internal/domain/model"
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

// --- stub implementations ---
type stubOssComponentRepo struct {
	createFn func(context.Context, *model.OssComponent) error
}

func (s *stubOssComponentRepo) Search(ctx context.Context, f domrepo.OssComponentFilter) ([]model.OssComponent, int, error) {
	return nil, 0, nil
}
func (s *stubOssComponentRepo) Create(ctx context.Context, c *model.OssComponent) error {
	if s.createFn != nil {
		return s.createFn(ctx, c)
	}
	return nil
}

type stubOssComponentLayerRepo struct {
	replaceFn func(context.Context, string, []string) error
}

func (s *stubOssComponentLayerRepo) ListByOssID(ctx context.Context, id string) ([]string, error) {
	return nil, nil
}
func (s *stubOssComponentLayerRepo) Replace(ctx context.Context, id string, layers []string) error {
	if s.replaceFn != nil {
		return s.replaceFn(ctx, id, layers)
	}
	return nil
}

type stubOssComponentTagRepo struct {
	replaceFn func(context.Context, string, []string) error
	listFn    func(context.Context, string) ([]model.Tag, error)
}

func (s *stubOssComponentTagRepo) Replace(ctx context.Context, id string, tagIDs []string) error {
	if s.replaceFn != nil {
		return s.replaceFn(ctx, id, tagIDs)
	}
	return nil
}
func (s *stubOssComponentTagRepo) ListByOssID(ctx context.Context, id string) ([]model.Tag, error) {
	if s.listFn != nil {
		return s.listFn(ctx, id)
	}
	return nil, nil
}

type stubOssVersionRepo struct {
	searchFn func(context.Context, domrepo.OssVersionFilter) ([]model.OssVersion, int, error)
	createFn func(context.Context, *model.OssVersion) error
	deleteFn func(context.Context, string) error
	getFn    func(context.Context, string) (*model.OssVersion, error)
	updateFn func(context.Context, *model.OssVersion) error
}

func (s *stubOssVersionRepo) Search(ctx context.Context, f domrepo.OssVersionFilter) ([]model.OssVersion, int, error) {
	if s.searchFn != nil {
		return s.searchFn(ctx, f)
	}
	return nil, 0, nil
}
func (s *stubOssVersionRepo) Create(ctx context.Context, v *model.OssVersion) error {
	if s.createFn != nil {
		return s.createFn(ctx, v)
	}
	return nil
}
func (s *stubOssVersionRepo) Delete(ctx context.Context, id string) error {
	if s.deleteFn != nil {
		return s.deleteFn(ctx, id)
	}
	return nil
}
func (s *stubOssVersionRepo) Get(ctx context.Context, id string) (*model.OssVersion, error) {
	if s.getFn != nil {
		return s.getFn(ctx, id)
	}
	return nil, nil
}
func (s *stubOssVersionRepo) Update(ctx context.Context, v *model.OssVersion) error {
	if s.updateFn != nil {
		return s.updateFn(ctx, v)
	}
	return nil
}

// --- tests ---
func TestToOssComponent_AllFields(t *testing.T) {
	now := time.Now()
	hp := "http://hp"
	repo := "http://repo"
	desc := "d"
	pl := "Go"
	role := "RUNTIME_REQUIRED"
	tagID := uuid.NewString()
	comp := model.OssComponent{ID: uuid.NewString(), Name: "Redis", NormalizedName: "redis", HomepageURL: &hp, RepositoryURL: &repo, Description: &desc, PrimaryLanguage: &pl, Layers: []string{"LIB"}, DefaultUsageRole: &role, Deprecated: false, CreatedAt: now, UpdatedAt: now, Tags: []model.Tag{{ID: tagID, Name: "db"}}}
	res := toOssComponent(comp)
	require.Equal(t, comp.ID, res.Id.String())
	require.NotNil(t, res.NormalizedName)
	require.Equal(t, "redis", *res.NormalizedName)
	require.NotNil(t, res.Layers)
	require.Equal(t, gen.Layer("LIB"), (*res.Layers)[0])
	require.NotNil(t, res.Tags)
	require.Equal(t, tagID, (*res.Tags)[0].Id.String())
	require.Equal(t, role, string(*res.DefaultUsageRole))
	require.Equal(t, &hp, res.HomepageUrl)
	require.Equal(t, &repo, res.RepositoryUrl)
	require.Equal(t, &desc, res.Description)
	require.Equal(t, &pl, res.PrimaryLanguage)
}

func TestToOssVersion_AllFields(t *testing.T) {
	now := time.Now()
	rel := now
	licRaw := "MIT"
	licConc := "MIT"
	purl := "pkg:"
	hash := "h"
	modDesc := "m"
	supp := "ORIGIN"
	fork := "f"
	version := model.OssVersion{ID: uuid.NewString(), OssID: uuid.NewString(), Version: "1", ReleaseDate: &rel, LicenseExpressionRaw: &licRaw, LicenseConcluded: &licConc, Purl: &purl, CpeList: []string{"c"}, HashSha256: &hash, Modified: true, ModificationDescription: &modDesc, ReviewStatus: "draft", LastReviewedAt: &now, ScopeStatus: "IN_SCOPE", SupplierType: &supp, ForkOriginURL: &fork, CreatedAt: now, UpdatedAt: now}
	res := toOssVersion(version)
	require.Equal(t, version.ID, res.Id.String())
	require.NotNil(t, res.ReleaseDate)
	require.NotNil(t, res.CpeList)
	require.Equal(t, "c", (*res.CpeList)[0])
	require.Equal(t, gen.ScopeStatus(version.ScopeStatus), res.ScopeStatus)
	require.Equal(t, gen.ReviewStatus(version.ReviewStatus), res.ReviewStatus)
	require.NotNil(t, res.HashSha256)
	require.Equal(t, hash, *res.HashSha256)
	require.Equal(t, gen.SupplierType(supp), *res.SupplierType)
	require.Equal(t, fork, *res.ForkOriginUrl)
}

func TestCreateOssComponent(t *testing.T) {
	var created *model.OssComponent
	tagID := uuid.NewString()
	compRepo := &stubOssComponentRepo{createFn: func(ctx context.Context, c *model.OssComponent) error { created = c; return nil }}
	layerCalled := false
	layerRepo := &stubOssComponentLayerRepo{replaceFn: func(ctx context.Context, id string, layers []string) error {
		layerCalled = true
		require.Equal(t, created.ID, id)
		require.Equal(t, []string{"LIB"}, layers)
		return nil
	}}
	tagCalled := false
	tagRepo := &stubOssComponentTagRepo{replaceFn: func(ctx context.Context, id string, ids []string) error {
		tagCalled = true
		require.Equal(t, created.ID, id)
		require.Equal(t, []string{tagID}, ids)
		return nil
	}, listFn: func(ctx context.Context, id string) ([]model.Tag, error) {
		return []model.Tag{{ID: tagID, Name: "db"}}, nil
	}}
	h := &Handler{OssComponentRepo: compRepo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)
	body := `{"name":"Redis","layers":["LIB"],"tagIds":["` + tagID + `"]}`
	req := httptest.NewRequest(http.MethodPost, "/oss", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.True(t, layerCalled)
	require.True(t, tagCalled)
	require.NotNil(t, created)
	var res gen.OssComponent
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, created.ID, res.Id.String())
}

func TestListOssVersions(t *testing.T) {
	ossID := uuid.NewString()
	now := time.Now()
	repo := &stubOssVersionRepo{searchFn: func(ctx context.Context, f domrepo.OssVersionFilter) ([]model.OssVersion, int, error) {
		require.Equal(t, ossID, f.OssID)
		return []model.OssVersion{{ID: uuid.NewString(), OssID: ossID, Version: "1", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}}, 1, nil
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodGet, "/oss/"+ossID+"/versions", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	var res gen.PagedResultOssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestCreateOssVersion(t *testing.T) {
	ossID := uuid.NewString()
	var created *model.OssVersion
	repo := &stubOssVersionRepo{createFn: func(ctx context.Context, v *model.OssVersion) error { created = v; return nil }}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"version":"1.0.0"}`
	req := httptest.NewRequest(http.MethodPost, "/oss/"+ossID+"/versions", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, created)
	require.Equal(t, ossID, created.OssID)
	var res gen.OssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, created.ID, res.Id.String())
}

func TestCreateOssVersion_WithOptions(t *testing.T) {
	ossID := uuid.NewString()
	var created *model.OssVersion
	repo := &stubOssVersionRepo{createFn: func(ctx context.Context, v *model.OssVersion) error {
		created = v
		return nil
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	now := time.Now().UTC().Truncate(time.Second)
	body := `{"version":"1.0.0","releaseDate":"` + now.Format("2006-01-02") + `","purl":"pkg:","cpeList":["c"],"modified":true,"supplierType":"ORIGIN"}`
	req := httptest.NewRequest(http.MethodPost, "/oss/"+ossID+"/versions", strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)
	require.NotNil(t, created)
	require.Equal(t, ossID, created.OssID)
	require.Equal(t, "pkg:", *created.Purl)
	require.Equal(t, []string{"c"}, created.CpeList)
	require.True(t, created.Modified)
	require.NotNil(t, created.ReleaseDate)
}

func TestCreateOssVersion_InvalidBody(t *testing.T) {
	ossID := uuid.NewString()
	repo := &stubOssVersionRepo{}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodPost, "/oss/"+ossID+"/versions", strings.NewReader("{"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteOssVersion(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	called := false
	repo := &stubOssVersionRepo{deleteFn: func(ctx context.Context, id string) error { called = true; require.Equal(t, vid, id); return nil }}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	req := httptest.NewRequest(http.MethodDelete, "/oss/"+ossID+"/versions/"+vid, nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
	require.True(t, called)
}

func TestUpdateOssVersion_NotFound(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	repo := &stubOssVersionRepo{getFn: func(ctx context.Context, id string) (*model.OssVersion, error) {
		return nil, sql.ErrNoRows
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"reviewStatus":"approved"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+ossID+"/versions/"+vid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNotFound, rec.Code)
}

func TestOssComponentStubMethods(t *testing.T) {
	h := &Handler{}
	e := echo.New()
	ctx := e.NewContext(httptest.NewRequest(http.MethodDelete, "/", nil), httptest.NewRecorder())
	require.NoError(t, h.DeprecateOssComponent(ctx, uuid.New()))
	ctx = e.NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	require.NoError(t, h.GetOssComponent(ctx, uuid.New()))
	ctx = e.NewContext(httptest.NewRequest(http.MethodPatch, "/", nil), httptest.NewRecorder())
	require.NoError(t, h.UpdateOssComponent(ctx, uuid.New()))
}

func TestUpdateOssVersion(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	now := time.Now()
	existing := model.OssVersion{ID: vid, OssID: ossID, Version: "1", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}
	var updated *model.OssVersion
	repo := &stubOssVersionRepo{getFn: func(ctx context.Context, id string) (*model.OssVersion, error) {
		require.Equal(t, vid, id)
		v := existing
		return &v, nil
	}, updateFn: func(ctx context.Context, v *model.OssVersion) error { updated = v; return nil }}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"reviewStatus":"approved"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+ossID+"/versions/"+vid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, updated)
	require.Equal(t, "approved", updated.ReviewStatus)
}

func TestUpdateOssVersion_WithFields(t *testing.T) {
	ossID := uuid.NewString()
	vid := uuid.NewString()
	now := time.Now()
	existing := model.OssVersion{ID: vid, OssID: ossID, Version: "1", ReviewStatus: "draft", ScopeStatus: "IN_SCOPE", CreatedAt: now, UpdatedAt: now}
	var updated *model.OssVersion
	repo := &stubOssVersionRepo{getFn: func(ctx context.Context, id string) (*model.OssVersion, error) {
		v := existing
		return &v, nil
	}, updateFn: func(ctx context.Context, v *model.OssVersion) error {
		updated = v
		return nil
	}}
	h := &Handler{OssVersionRepo: repo}
	e := setupEcho(h)
	body := `{"licenseExpressionRaw":"MIT","scopeStatus":"OUT_SCOPE"}`
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+ossID+"/versions/"+vid, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.NotNil(t, updated)
	require.Equal(t, "MIT", *updated.LicenseExpressionRaw)
	require.Equal(t, "OUT_SCOPE", updated.ScopeStatus)
}
