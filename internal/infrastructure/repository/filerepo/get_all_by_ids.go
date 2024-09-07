package filerepo

import (
	"context"
	"encoding/json"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/file/contenttype"
	"github.com/puny-activity/file-service/internal/entity/file/path"
	"github.com/puny-activity/file-service/pkg/queryer"
	"github.com/puny-activity/file-service/pkg/werr"
)

type getAllByIDsEntity struct {
	ID          string          `db:"id"`
	RootName    string          `db:"root_name"`
	Path        string          `db:"path"`
	Name        string          `db:"name"`
	ContentType string          `db:"content_type"`
	Size        int64           `db:"size"`
	Metadata    json.RawMessage `db:"metadata"`
	MD5         string          `db:"md5"`
}

func (r *Repository) GetAllByIDs(ctx context.Context, fileIDs []file.ID) ([]file.File, error) {
	return r.getAllByIDs(ctx, r.db, fileIDs)
}

func (r *Repository) GetAllByIDsTx(ctx context.Context, tx *sqlx.Tx, fileIDs []file.ID) ([]file.File, error) {
	return r.getAllByIDs(ctx, tx, fileIDs)
}

func (r *Repository) getAllByIDs(ctx context.Context, queryer queryer.Queryer, fileIDs []file.ID) ([]file.File, error) {
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
         JOIN roots r on f.root_id = r.id
WHERE f.id = ANY($1)
`

	fileIDsString := make([]string, len(fileIDs))
	for i, id := range fileIDs {
		fileIDsString[i] = id.String()
	}

	filesRepo := make([]getAllByIDsEntity, 0)
	err := queryer.SelectContext(ctx, &filesRepo, query, fileIDsString)
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
			Path:        path.New(filesRepo[i].RootName, filesRepo[i].Path),
			Name:        filesRepo[i].Name,
			ContentType: contentType,
			Size:        filesRepo[i].Size,
			Metadata:    filesRepo[i].Metadata,
			MD5:         filesRepo[i].MD5,
		}
	}

	return files, nil
}
