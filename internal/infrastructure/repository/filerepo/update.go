package filerepo

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/pkg/queryer"
)

type updateEntity struct {
	ID          uuid.UUID       `db:"id"`
	Path        string          `db:"path"`
	Name        string          `db:"name"`
	ContentType string          `db:"content_type"`
	Size        int64           `db:"size"`
	Metadata    json.RawMessage `db:"metadata"`
	MD5         string          `db:"md5"`
}

func (r *Repository) Update(ctx context.Context, fileToCreate file.File) error {
	return r.update(ctx, r.db, fileToCreate)
}

func (r *Repository) UpdateTx(ctx context.Context, tx *sqlx.Tx, fileToCreate file.File) error {
	return r.update(ctx, tx, fileToCreate)
}

func (r *Repository) update(ctx context.Context, queryer queryer.Queryer, fileToCreate file.File) error {
	query := `
UPDATE files f
SET path         = :path,
    name         = :name,
    content_type = :content_type,
    size         = :size,
    metadata     = :metadata,
    md5          = :md5
WHERE f.id = :id
`

	parameter := updateEntity{
		ID:          uuid.UUID(*fileToCreate.ID),
		Path:        fileToCreate.Path,
		Name:        fileToCreate.Name,
		ContentType: fileToCreate.ContentType.String(),
		Size:        fileToCreate.Size,
		Metadata:    fileToCreate.Metadata,
		MD5:         fileToCreate.MD5,
	}

	_, err := queryer.NamedExecContext(ctx, query, parameter)
	if err != nil {
		return err
	}

	return nil
}
