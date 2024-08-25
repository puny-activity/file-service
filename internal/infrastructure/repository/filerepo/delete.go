package filerepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/pkg/queryer"
)

func (r *Repository) Delete(ctx context.Context, fileID file.ID) error {
	return r.delete(ctx, r.db, fileID)
}

func (r *Repository) DeleteTx(ctx context.Context, tx *sqlx.Tx, fileID file.ID) error {
	return r.delete(ctx, tx, fileID)
}

func (r *Repository) delete(ctx context.Context, queryer queryer.Queryer, fileID file.ID) error {
	query := `
DELETE
FROM files
WHERE id = $1
`
	_, err := queryer.ExecContext(ctx, query, uuid.UUID(fileID))
	if err != nil {
		return err
	}

	return nil
}
