package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

// ProjectUsageRepository implements domrepo.ProjectUsageRepository.
type ProjectUsageRepository struct {
	DB *sql.DB
}

var _ domrepo.ProjectUsageRepository = (*ProjectUsageRepository)(nil)

// Search returns project usages matching filter.
func (r *ProjectUsageRepository) Search(ctx context.Context, f domrepo.ProjectUsageFilter) ([]model.ProjectUsage, int, error) {
	var args []any
	wheres := []string{"project_id = ?"}
	args = append(args, f.ProjectID)
	if f.ScopeStatus != "" {
		wheres = append(wheres, "scope_status = ?")
		args = append(args, f.ScopeStatus)
	}
	if f.UsageRole != "" {
		wheres = append(wheres, "usage_role = ?")
		args = append(args, f.UsageRole)
	}
	if f.Direct != nil {
		wheres = append(wheres, "direct_dependency = ?")
		args = append(args, *f.Direct)
	}
	whereSQL := "WHERE " + strings.Join(wheres, " AND ")

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM project_usages %s", whereSQL)
	row := r.DB.QueryRowContext(ctx, countQuery, args...)
	var total int
	if err := row.Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (f.Page - 1) * f.Size
	listQuery := fmt.Sprintf(`SELECT id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by FROM project_usages %s ORDER BY added_at DESC LIMIT ? OFFSET ?`, whereSQL)
	argsWithLimit := append(args, f.Size, offset)
	rows, err := r.DB.QueryContext(ctx, listQuery, argsWithLimit...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var usages []model.ProjectUsage
	for rows.Next() {
		var u model.ProjectUsage
		var note, evalBy sql.NullString
		var evalAt sql.NullTime
		if err := rows.Scan(&u.ID, &u.ProjectID, &u.OssID, &u.OssVersionID, &u.UsageRole, &u.ScopeStatus, &note, &u.DirectDependency, &u.AddedAt, &evalAt, &evalBy); err != nil {
			return nil, 0, err
		}
		if note.Valid {
			u.InclusionNote = &note.String
		}
		if evalAt.Valid {
			u.EvaluatedAt = &evalAt.Time
		}
		if evalBy.Valid {
			u.EvaluatedBy = &evalBy.String
		}
		usages = append(usages, u)
	}
	return usages, total, rows.Err()
}

// Create inserts a new project usage.
func (r *ProjectUsageRepository) Create(ctx context.Context, u *model.ProjectUsage) error {
	query := `INSERT INTO project_usages (id, project_id, oss_id, oss_version_id, usage_role, scope_status, inclusion_note, direct_dependency, added_at, evaluated_at, evaluated_by) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.DB.ExecContext(ctx, query, u.ID, u.ProjectID, u.OssID, u.OssVersionID, u.UsageRole, u.ScopeStatus, u.InclusionNote, u.DirectDependency, u.AddedAt, u.EvaluatedAt, u.EvaluatedBy)
	return err
}

// Update updates an existing usage.
func (r *ProjectUsageRepository) Update(ctx context.Context, u *model.ProjectUsage) error {
	query := `UPDATE project_usages SET oss_version_id = ?, usage_role = ?, direct_dependency = ?, inclusion_note = ?, scope_status = ?, evaluated_at = ?, evaluated_by = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, u.OssVersionID, u.UsageRole, u.DirectDependency, u.InclusionNote, u.ScopeStatus, u.EvaluatedAt, u.EvaluatedBy, u.ID)
	return err
}

// Delete removes a usage by ID.
func (r *ProjectUsageRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM project_usages WHERE id = ?`, id)
	return err
}

// UpdateScope updates only scope related fields.
func (r *ProjectUsageRepository) UpdateScope(ctx context.Context, id string, scopeStatus string, inclusionNote *string, evaluatedAt time.Time, evaluatedBy *string) error {
	query := `UPDATE project_usages SET scope_status = ?, inclusion_note = ?, evaluated_at = ?, evaluated_by = ? WHERE id = ?`
	_, err := r.DB.ExecContext(ctx, query, scopeStatus, inclusionNote, evaluatedAt, evaluatedBy, id)
	return err
}
