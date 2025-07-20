package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

// OssComponentTagRepository implements domrepo.OssComponentTagRepository.
type OssComponentTagRepository struct {
	DB *sql.DB
}

var _ domrepo.OssComponentTagRepository = (*OssComponentTagRepository)(nil)

// ListByOssID returns tags associated with the given component ordered by created_at.
func (r *OssComponentTagRepository) ListByOssID(ctx context.Context, ossID string) ([]model.Tag, error) {
	rows, err := r.DB.QueryContext(ctx,
		`SELECT tg.id, tg.name, tg.created_at FROM tags tg
         JOIN oss_component_tags ct ON ct.tag_id = tg.id
         WHERE ct.oss_id = ? ORDER BY tg.created_at DESC`, ossID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []model.Tag
	for rows.Next() {
		var t model.Tag
		var created time.Time
		if err := rows.Scan(&t.ID, &t.Name, &created); err != nil {
			return nil, err
		}
		t.CreatedAt = &created
		tags = append(tags, t)
	}
	return tags, rows.Err()
}

// Replace replaces tags for a component with provided IDs.
func (r *OssComponentTagRepository) Replace(ctx context.Context, ossID string, tagIDs []string) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM oss_component_tags WHERE oss_id = ?`, ossID); err != nil {
		tx.Rollback()
		return err
	}
	for _, id := range tagIDs {
		if _, err := tx.ExecContext(ctx, `INSERT INTO oss_component_tags (oss_id, tag_id) VALUES (?, ?)`, ossID, id); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
