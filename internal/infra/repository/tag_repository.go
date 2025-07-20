package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/ramsesyok/oss-catalog/internal/domain/model"
	domrepo "github.com/ramsesyok/oss-catalog/internal/domain/repository"
)

// TagRepository implements domrepo.TagRepository.
type TagRepository struct {
	DB *sql.DB
}

var _ domrepo.TagRepository = (*TagRepository)(nil)

// List returns all tags ordered by created_at descending.
func (r *TagRepository) List(ctx context.Context) ([]model.Tag, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, name, created_at FROM tags ORDER BY created_at DESC`)
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

// Create inserts a new tag.
func (r *TagRepository) Create(ctx context.Context, t *model.Tag) error {
	_, err := r.DB.ExecContext(ctx,
		`INSERT INTO tags (id, name, created_at) VALUES (?, ?, ?)`,
		t.ID, t.Name, t.CreatedAt,
	)
	return err
}

// Delete removes a tag by ID.
func (r *TagRepository) Delete(ctx context.Context, id string) error {
	_, err := r.DB.ExecContext(ctx, `DELETE FROM tags WHERE id = ?`, id)
	return err
}
