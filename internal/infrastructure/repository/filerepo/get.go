package filerepo

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/file/filecontenttype"
	"github.com/puny-activity/file-service/pkg/queryer"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/puny-activity/file-service/pkg/werr"
)

type getEntity struct {
	ID          uuid.UUID       `db:"id"`
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
SELECT id,
       path,
       name,
       content_type,
       size,
       metadata,
       md5
FROM files f
WHERE f.id = $1
`

	var fileRepo getEntity
	err := queryer.GetContext(ctx, &fileRepo, query, uuid.UUID(fileID))
	if err != nil {
		return file.File{}, err
	}

	contentType, err := filecontenttype.New(fileRepo.ContentType)
	if err != nil {
		return file.File{}, werr.WrapSE("failed to parse content type", err)
	}

	file := file.File{
		ID:          util.ToPointer(file.ID(fileRepo.ID)),
		Path:        fileRepo.Path,
		Name:        fileRepo.Name,
		ContentType: contentType,
		Size:        fileRepo.Size,
		Metadata:    fileRepo.Metadata,
		MD5:         fileRepo.MD5,
	}

	return file, nil
}
