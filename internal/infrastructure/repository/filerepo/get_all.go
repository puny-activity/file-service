package filerepo

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/file/filecontenttype"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/queryer"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/puny-activity/file-service/pkg/werr"
)

type getAllEntity struct {
	ID          uuid.UUID       `db:"id"`
	Path        string          `db:"path"`
	Name        string          `db:"name"`
	ContentType string          `db:"content_type"`
	Size        int64           `db:"size"`
	Metadata    json.RawMessage `db:"metadata"`
	MD5         string          `db:"md5"`
}

func (r *Repository) GetAll(ctx context.Context, rootID root.ID) ([]file.File, error) {
	return r.getAll(ctx, r.db, rootID)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx, rootID root.ID) ([]file.File, error) {
	return r.getAll(ctx, tx, rootID)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer, rootID root.ID) ([]file.File, error) {
	query := `
SELECT id,
       path,
       name,
       content_type,
       size,
       metadata,
       md5
FROM files f
WHERE f.root_id = $1
`

	filesRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &filesRepo, query, rootID.String())
	if err != nil {
		return nil, err
	}

	files := make([]file.File, len(filesRepo))
	for i := range filesRepo {
		contentType, err := filecontenttype.New(filesRepo[i].ContentType)
		if err != nil {
			return nil, werr.WrapSE("failed to parse content type", err)
		}

		files[i] = file.File{
			ID:          util.ToPointer(file.ID(filesRepo[i].ID)),
			Path:        filesRepo[i].Path,
			Name:        filesRepo[i].Name,
			ContentType: contentType,
			Size:        filesRepo[i].Size,
			Metadata:    filesRepo[i].Metadata,
			MD5:         filesRepo[i].MD5,
		}
	}

	return files, nil
}
