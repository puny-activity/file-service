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

type getEntity struct {
	ID          string          `db:"id"`
	RootName    string          `db:"root_name"`
	Path        string          `db:"path"`
	Name        string          `db:"name"`
	ContentType string          `db:"content_type"`
	Size        int64           `db:"size"`
	Metadata    json.RawMessage `db:"metadata"`
	MD5         string          `db:"md5"`
}

func (r *Repository) Get(ctx context.Context, fileID file.ID) (file.File, error) {
	return r.get(ctx, r.db, fileID)
}

func (r *Repository) GetTx(ctx context.Context, tx *sqlx.Tx, fileID file.ID) (file.File, error) {
	return r.get(ctx, tx, fileID)
}

func (r *Repository) get(ctx context.Context, queryer queryer.Queryer, fileID file.ID) (file.File, error) {
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
JOIN roots r ON r.id = f.root_id
WHERE f.id = $1
`

	var fileRepo getEntity
	err := queryer.GetContext(ctx, &fileRepo, query, fileID.String())
	if err != nil {
		return file.File{}, err
	}

	contentType, err := contenttype.New(fileRepo.ContentType)
	if err != nil {
		return file.File{}, werr.WrapSE("failed to parse content type", err)
	}

	file := file.File{
		ID:          &fileID,
		Path:        path.New(fileRepo.RootName, fileRepo.Path),
		Name:        fileRepo.Name,
		ContentType: contentType,
		Size:        fileRepo.Size,
		Metadata:    fileRepo.Metadata,
		MD5:         fileRepo.MD5,
	}

	return file, nil
}
