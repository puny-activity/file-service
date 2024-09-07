package filerepo

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/queryer"
)

type createEntity struct {
	ID          string          `db:"id"`
	RootID      string          `db:"root_id"`
	Path        string          `db:"path"`
	Name        string          `db:"name"`
	ContentType string          `db:"content_type"`
	Size        int64           `db:"size"`
	Metadata    json.RawMessage `db:"metadata"`
	MD5         string          `db:"md5"`
}

func (r *Repository) Create(ctx context.Context, rootID root.ID, fileToCreate file.File) error {
	return r.create(ctx, r.db, rootID, fileToCreate)
}

func (r *Repository) CreateTx(ctx context.Context, tx *sqlx.Tx, rootID root.ID, fileToCreate file.File) error {
	return r.create(ctx, tx, rootID, fileToCreate)
}

func (r *Repository) create(ctx context.Context, queryer queryer.Queryer, rootID root.ID, fileToCreate file.File) error {
	query := `
INSERT INTO files(id, root_id, path, name, content_type, size, metadata, md5)
VALUES (:id, :root_id, :path, :name, :content_type, :size, :metadata, :md5)
`

	parameter := createEntity{
		ID:          (*fileToCreate.ID).String(),
		RootID:      rootID.String(),
		Path:        fileToCreate.Path.RelativePath(),
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
