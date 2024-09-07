package filerepo

import (
	"context"
	"encoding/json"
	sqlx "github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/file/contenttype"
	"github.com/puny-activity/file-service/internal/entity/file/path"
	"github.com/puny-activity/file-service/pkg/queryer"
	"github.com/puny-activity/file-service/pkg/werr"
)

type getAllEntity struct {
	ID          string          `db:"id"`
	RootID      string          `db:"root_id"`
	Path        string          `db:"path"`
	Name        string          `db:"name"`
	ContentType string          `db:"content_type"`
	Size        int64           `db:"size"`
	Metadata    json.RawMessage `db:"metadata"`
	MD5         string          `db:"md5"`
}

func (r *Repository) GetAll(ctx context.Context) ([]file.File, error) {
	return r.getAll(ctx, r.db)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]file.File, error) {
	return r.getAll(ctx, tx)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer) ([]file.File, error) {
	query := `
SELECT f.id,
       r.name AS root_name,
       path,
       f.name,
       content_type,
       size,
       metadata,
       md5
FROM files f
         JOIN roots r on r.id = f.root_id
`

	filesRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &filesRepo, query)
	if err != nil {
		return nil, err
	}

	files := make([]file.File, len(filesRepo))
	for i := range filesRepo {
		fileID, err := file.ParseID(filesRepo[i].ID)
		if err != nil {
			return nil, werr.WrapSE("failed to parse file id", err)
		}

		contentType, err := contenttype.New(filesRepo[i].ContentType)
		if err != nil {
			return nil, werr.WrapSE("failed to parse content type", err)
		}

		files[i] = file.File{
			ID:          &fileID,
			Path:        path.New(filesRepo[i].Name, filesRepo[i].Path),
			Name:        filesRepo[i].Name,
			ContentType: contentType,
			Size:        filesRepo[i].Size,
			Metadata:    filesRepo[i].Metadata,
			MD5:         filesRepo[i].MD5,
		}
	}

	return files, nil
}
