package filehistoryrepo

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/filehistory"
	"github.com/puny-activity/file-service/internal/entity/filehistory/actiontype"
	"github.com/puny-activity/file-service/pkg/queryer"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/puny-activity/file-service/pkg/werr"
)

type getSinceEntity struct {
	ID          uuid.UUID `db:"id"`
	FileID      uuid.UUID `db:"file_id"`
	Action      string    `db:"action_type"`
	PerformedAt string    `db:"performed_at"`
}

func (r *Repository) GetSince(ctx context.Context, updatedSince carbon.Carbon) ([]filehistory.Row, error) {
	return r.getSince(ctx, r.db, updatedSince)
}

func (r *Repository) GetSinceTx(ctx context.Context, tx *sqlx.Tx, updatedSince carbon.Carbon) ([]filehistory.Row, error) {
	return r.getSince(ctx, tx, updatedSince)
}

func (r *Repository) getSince(ctx context.Context, queryer queryer.Queryer, updatedSince carbon.Carbon) ([]filehistory.Row, error) {
	query := `
SELECT fh.id,
       fh.file_id,
       fh.action_type,
       fh.performed_at
FROM file_history fh
`

	rowsRepo := make([]getSinceEntity, 0)
	err := queryer.SelectContext(ctx, &rowsRepo, query)
	if err != nil {
		return nil, err
	}

	rows := make([]filehistory.Row, len(rowsRepo))
	for i := range rowsRepo {
		action, err := actiontype.New(rowsRepo[i].Action)
		if err != nil {
			return nil, werr.WrapSE("failed to parse content type", err)
		}

		performedAt := carbon.Parse(rowsRepo[i].PerformedAt)
		if performedAt.Error != nil {
			return nil, werr.WrapSE("failed to parse content time", performedAt.Error)
		}

		rows[i] = filehistory.Row{
			ID:          util.ToPointer(filehistory.ID(rowsRepo[i].ID)),
			FileID:      file.ID(rowsRepo[i].FileID),
			Action:      action,
			PerformedAt: performedAt,
		}
	}

	return rows, nil
}
