package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/stretchr/testify/require"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"

	gen "github.com/ramsesyok/oss-catalog/internal/api/gen"
	infrarepo "github.com/ramsesyok/oss-catalog/internal/infra/repository"
)

func setupEcho(h *Handler) *echo.Echo {
	e := echo.New()
	gen.RegisterHandlers(e, h)
	return e
}

func TestListOssComponents(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	compRepo := &infrarepo.OssComponentRepository{DB: db}
	layerRepo := &infrarepo.OssComponentLayerRepository{DB: db}
	tagRepo := &infrarepo.OssComponentTagRepository{DB: db}

	h := &Handler{OssComponentRepo: compRepo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)

	id := uuid.NewString()
	now := time.Now()
	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_components oc WHERE normalized_name LIKE ?")
	mock.ExpectQuery(countQuery).WithArgs("%redis%").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT oc.id, oc.name, oc.normalized_name, oc.homepage_url, oc.repository_url, oc.description, oc.primary_language, oc.default_usage_role, oc.deprecated, oc.created_at, oc.updated_at FROM oss_components oc WHERE normalized_name LIKE ? ORDER BY oc.created_at DESC LIMIT ? OFFSET ?")
	mock.ExpectQuery(listQuery).WithArgs("%redis%", 50, 0).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "normalized_name", "homepage_url", "repository_url", "description", "primary_language", "default_usage_role", "deprecated", "created_at", "updated_at"}).
			AddRow(id, "Redis", "redis", nil, nil, nil, nil, nil, false, now, now))

	mock.ExpectQuery(regexp.QuoteMeta("SELECT layer FROM oss_component_layers WHERE oss_id = ? ORDER BY layer")).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"layer"}))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT tg.id, tg.name, tg.created_at FROM tags tg JOIN oss_component_tags ct ON ct.tag_id = tg.id WHERE ct.oss_id = ? ORDER BY tg.created_at DESC")).
		WithArgs(id).WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}))

	req := httptest.NewRequest(http.MethodGet, "/oss?name=redis", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.PagedResultOssComponent
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestGetOssVersion_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	vid := uuid.New()
	query := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at FROM oss_versions WHERE id = ?")
	mock.ExpectQuery(query).WithArgs(vid.String()).WillReturnError(sql.ErrNoRows)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+uuid.New().String()+"/versions/"+vid.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNotFound, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOssVersion_OK(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	vid := uuid.New()
	oid := uuid.New()
	now := time.Now()
	query := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at FROM oss_versions WHERE id = ?")
	mockRows := sqlmock.NewRows([]string{"id", "oss_id", "version", "release_date", "license_expression_raw", "license_concluded", "purl", "cpe_list", "hash_sha256", "modified", "modification_description", "review_status", "last_reviewed_at", "scope_status", "supplier_type", "fork_origin_url", "created_at", "updated_at"}).
		AddRow(vid.String(), oid.String(), "1.0.0", now, nil, nil, nil, pq.StringArray{}, nil, false, nil, "draft", nil, "IN_SCOPE", nil, nil, now, now)
	mock.ExpectQuery(query).WithArgs(vid.String()).WillReturnRows(mockRows)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+oid.String()+"/versions/"+vid.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.OssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.Equal(t, vid, res.Id)
}
func TestCreateOssComponent(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	compRepo := &infrarepo.OssComponentRepository{DB: db}
	layerRepo := &infrarepo.OssComponentLayerRepository{DB: db}
	tagRepo := &infrarepo.OssComponentTagRepository{DB: db}
	h := &Handler{OssComponentRepo: compRepo, OssComponentLayerRepo: layerRepo, OssComponentTagRepo: tagRepo}
	e := setupEcho(h)

	tagID := uuid.New()

	insQuery := regexp.QuoteMeta("INSERT INTO oss_components (id, name, normalized_name, homepage_url, repository_url, description, primary_language, default_usage_role, deprecated, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(insQuery).
		WithArgs(sqlmock.AnyArg(), "Redis", "redis", nil, nil, nil, nil, nil, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectBegin()
	delLayer := regexp.QuoteMeta("DELETE FROM oss_component_layers WHERE oss_id = ?")
	mock.ExpectExec(delLayer).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	insLayer := regexp.QuoteMeta("INSERT INTO oss_component_layers (oss_id, layer) VALUES (?, ?)")
	mock.ExpectExec(insLayer).WithArgs(sqlmock.AnyArg(), "LIB").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	delTag := regexp.QuoteMeta("DELETE FROM oss_component_tags WHERE oss_id = ?")
	mock.ExpectExec(delTag).WithArgs(sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	insTag := regexp.QuoteMeta("INSERT INTO oss_component_tags (oss_id, tag_id) VALUES (?, ?)")
	mock.ExpectExec(insTag).WithArgs(sqlmock.AnyArg(), tagID.String()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	listTag := regexp.QuoteMeta("SELECT tg.id, tg.name, tg.created_at FROM tags tg JOIN oss_component_tags ct ON ct.tag_id = tg.id WHERE ct.oss_id = ? ORDER BY tg.created_at DESC")
	mock.ExpectQuery(listTag).
		WithArgs(sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "created_at"}).
			AddRow(tagID.String(), "db", time.Now()))

	body, _ := json.Marshal(gen.OssComponentCreateRequest{
		Name:   "Redis",
		Layers: &[]gen.Layer{gen.LIB},
		TagIds: &[]openapi_types.UUID{tagID},
	})
	req := httptest.NewRequest(http.MethodPost, "/oss", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListOssVersions(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	oid := uuid.New()
	countQuery := regexp.QuoteMeta("SELECT COUNT(*) FROM oss_versions WHERE oss_id = ?")
	mock.ExpectQuery(countQuery).WithArgs(oid.String()).WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	listQuery := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at FROM oss_versions WHERE oss_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?")
	now := time.Now()
	rows := sqlmock.NewRows([]string{"id", "oss_id", "version", "release_date", "license_expression_raw", "license_concluded", "purl", "cpe_list", "hash_sha256", "modified", "modification_description", "review_status", "last_reviewed_at", "scope_status", "supplier_type", "fork_origin_url", "created_at", "updated_at"}).
		AddRow(uuid.NewString(), oid.String(), "1.0.0", now, nil, nil, nil, pq.StringArray{}, nil, false, nil, "draft", nil, "IN_SCOPE", nil, nil, now, now)
	mock.ExpectQuery(listQuery).WithArgs(oid.String(), 50, 0).WillReturnRows(rows)

	req := httptest.NewRequest(http.MethodGet, "/oss/"+oid.String()+"/versions", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
	var res gen.PagedResultOssVersion
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &res))
	require.NotNil(t, res.Items)
	require.Len(t, *res.Items, 1)
}

func TestCreateOssVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	oid := uuid.New()
	insertQuery := regexp.QuoteMeta("INSERT INTO oss_versions (id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(insertQuery).
		WithArgs(sqlmock.AnyArg(), oid.String(), "1.0.0", sqlmock.AnyArg(), nil, nil, nil, sqlmock.AnyArg(), nil, false, nil, "draft", nil, "IN_SCOPE", nil, nil, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	now := time.Now()
	body, _ := json.Marshal(gen.OssVersionCreateRequest{Version: "1.0.0", ReleaseDate: &openapi_types.Date{Time: now}})
	req := httptest.NewRequest(http.MethodPost, "/oss/"+oid.String()+"/versions", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOssVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	oid := uuid.New()
	vid := uuid.New()
	now := time.Now()
	getQuery := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at FROM oss_versions WHERE id = ?")
	mockRows := sqlmock.NewRows([]string{"id", "oss_id", "version", "release_date", "license_expression_raw", "license_concluded", "purl", "cpe_list", "hash_sha256", "modified", "modification_description", "review_status", "last_reviewed_at", "scope_status", "supplier_type", "fork_origin_url", "created_at", "updated_at"}).
		AddRow(vid.String(), oid.String(), "1.0.0", now, nil, nil, nil, pq.StringArray{}, nil, false, nil, "draft", nil, "IN_SCOPE", nil, nil, now, now)
	mock.ExpectQuery(getQuery).WithArgs(vid.String()).WillReturnRows(mockRows)

	updateQuery := regexp.QuoteMeta("UPDATE oss_versions SET release_date = ?, license_expression_raw = ?, license_concluded = ?, purl = ?, cpe_list = ?, hash_sha256 = ?, modified = ?, modification_description = ?, review_status = ?, last_reviewed_at = ?, scope_status = ?, supplier_type = ?, fork_origin_url = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(updateQuery).
		WithArgs(sqlmock.AnyArg(), nil, nil, nil, sqlmock.AnyArg(), nil, true, nil, "verified", sqlmock.AnyArg(), "IN_SCOPE", nil, nil, sqlmock.AnyArg(), vid.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	status := gen.Verified
	body, _ := json.Marshal(gen.OssVersionUpdateRequest{Modified: boolPtr(true), ReviewStatus: &status})
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+oid.String()+"/versions/"+vid.String(), bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteOssVersion(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	vid := uuid.New()
	delQuery := regexp.QuoteMeta("DELETE FROM oss_versions WHERE id = ?")
	mock.ExpectExec(delQuery).WithArgs(vid.String()).WillReturnResult(sqlmock.NewResult(1, 1))

	req := httptest.NewRequest(http.MethodDelete, "/oss/"+uuid.New().String()+"/versions/"+vid.String(), nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusNoContent, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateOssVersion_AllFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	oid := uuid.New()
	insertQuery := regexp.QuoteMeta("INSERT INTO oss_versions (id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	mock.ExpectExec(insertQuery).
		WithArgs(sqlmock.AnyArg(), oid.String(), "2.0.0", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), true, sqlmock.AnyArg(), "draft", sqlmock.AnyArg(), "IN_SCOPE", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	now := time.Now()
	mod := true
	reqBody, _ := json.Marshal(gen.OssVersionCreateRequest{
		Version:                 "2.0.0",
		ReleaseDate:             &openapi_types.Date{Time: now},
		LicenseExpressionRaw:    strPtr("Apache-2.0"),
		Purl:                    strPtr("pkg:maven/org.example/app@2.0.0"),
		CpeList:                 &[]string{"cpe:/a"},
		HashSha256:              strPtr("deadbeef"),
		Modified:                &mod,
		ModificationDescription: strPtr("mod"),
		SupplierType:            func() *gen.SupplierType { v := gen.INTERNALFORK; return &v }(),
		ForkOriginUrl:           strPtr("https://upstream"),
	})
	req := httptest.NewRequest(http.MethodPost, "/oss/"+oid.String()+"/versions", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateOssVersion_AllFields(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	verRepo := &infrarepo.OssVersionRepository{DB: db}
	h := &Handler{OssVersionRepo: verRepo}
	e := setupEcho(h)

	oid := uuid.New()
	vid := uuid.New()
	now := time.Now()
	getQuery := regexp.QuoteMeta("SELECT id, oss_id, version, release_date, license_expression_raw, license_concluded, purl, cpe_list, hash_sha256, modified, modification_description, review_status, last_reviewed_at, scope_status, supplier_type, fork_origin_url, created_at, updated_at FROM oss_versions WHERE id = ?")
	mockRows := sqlmock.NewRows([]string{"id", "oss_id", "version", "release_date", "license_expression_raw", "license_concluded", "purl", "cpe_list", "hash_sha256", "modified", "modification_description", "review_status", "last_reviewed_at", "scope_status", "supplier_type", "fork_origin_url", "created_at", "updated_at"}).
		AddRow(vid.String(), oid.String(), "2.0.0", now, nil, nil, nil, pq.StringArray{}, nil, false, nil, "draft", nil, "IN_SCOPE", nil, nil, now, now)
	mock.ExpectQuery(getQuery).WithArgs(vid.String()).WillReturnRows(mockRows)

	updateQuery := regexp.QuoteMeta("UPDATE oss_versions SET release_date = ?, license_expression_raw = ?, license_concluded = ?, purl = ?, cpe_list = ?, hash_sha256 = ?, modified = ?, modification_description = ?, review_status = ?, last_reviewed_at = ?, scope_status = ?, supplier_type = ?, fork_origin_url = ?, updated_at = ? WHERE id = ?")
	mock.ExpectExec(updateQuery).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), true, sqlmock.AnyArg(), "verified", sqlmock.AnyArg(), "OUT_SCOPE", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), vid.String()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mod := true
	status := gen.Verified
	scope := gen.OUTSCOPE
	reqBody, _ := json.Marshal(gen.OssVersionUpdateRequest{
		ReleaseDate:             &openapi_types.Date{Time: now},
		LicenseExpressionRaw:    strPtr("Apache-2.0"),
		LicenseConcluded:        strPtr("Apache-2.0"),
		Purl:                    strPtr("pkg:maven/org.example/app@2.0.0"),
		CpeList:                 &[]string{"cpe:/a"},
		HashSha256:              strPtr("deadbeef"),
		Modified:                &mod,
		ModificationDescription: strPtr("mod"),
		ReviewStatus:            &status,
		ScopeStatus:             &scope,
		SupplierType:            func() *gen.SupplierType { v := gen.INTERNALFORK; return &v }(),
		ForkOriginUrl:           strPtr("https://upstream"),
	})
	req := httptest.NewRequest(http.MethodPatch, "/oss/"+oid.String()+"/versions/"+vid.String(), bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.NoError(t, mock.ExpectationsWereMet())
}

func strPtr(s string) *string { return &s }

func boolPtr(b bool) *bool { return &b }

func TestToOssComponent(t *testing.T) {
	id := uuid.NewString()
	now := time.Now()
	homepage := "https://example.com"
	repo := "https://github.com/example"
	desc := "desc"
	lang := "Go"
	role := "STATIC_LINK"
	tagID := uuid.NewString()
	comp := model.OssComponent{
		ID:               id,
		Name:             "Example",
		NormalizedName:   "example",
		HomepageURL:      &homepage,
		RepositoryURL:    &repo,
		Description:      &desc,
		PrimaryLanguage:  &lang,
		Layers:           []string{"LIB", "DB"},
		DefaultUsageRole: &role,
		Tags:             []model.Tag{{ID: tagID, Name: "db", CreatedAt: &now}},
		Deprecated:       false,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	res := toOssComponent(comp)
	require.Equal(t, uuid.MustParse(id), res.Id)
	require.NotNil(t, res.Layers)
	require.NotNil(t, res.Tags)
	require.Len(t, *res.Tags, 1)
}

func TestToOssVersion(t *testing.T) {
	id := uuid.NewString()
	oid := uuid.NewString()
	now := time.Now()
	license := "MIT"
	conc := "MIT"
	purl := "pkg:generic/oss@1.0.0"
	hash := "deadbeef"
	modDesc := "changed"
	supplier := "INTERNAL_FORK"
	fork := "https://origin"
	v := model.OssVersion{
		ID:                      id,
		OssID:                   oid,
		Version:                 "1.0.0",
		ReleaseDate:             &now,
		LicenseExpressionRaw:    &license,
		LicenseConcluded:        &conc,
		Purl:                    &purl,
		CpeList:                 []string{"cpe:/a"},
		HashSha256:              &hash,
		ModificationDescription: &modDesc,
		LastReviewedAt:          &now,
		Modified:                true,
		ReviewStatus:            "verified",
		ScopeStatus:             "IN_SCOPE",
		SupplierType:            &supplier,
		ForkOriginURL:           &fork,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
	res := toOssVersion(v)
	require.Equal(t, uuid.MustParse(id), res.Id)
	require.NotNil(t, res.ReleaseDate)
	require.NotNil(t, res.LicenseExpressionRaw)
	require.NotNil(t, res.LicenseConcluded)
	require.NotNil(t, res.Purl)
	require.NotNil(t, res.CpeList)
	require.NotNil(t, res.HashSha256)
	require.NotNil(t, res.ModificationDescription)
	require.NotNil(t, res.LastReviewedAt)
	require.NotNil(t, res.SupplierType)
	require.NotNil(t, res.ForkOriginUrl)
}

func TestOssComponentNoopHandlers(t *testing.T) {
	h := &Handler{}
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	require.NoError(t, h.DeprecateOssComponent(c, uuid.New()))
	require.NoError(t, h.GetOssComponent(c, uuid.New()))
	require.NoError(t, h.UpdateOssComponent(c, uuid.New()))
}
