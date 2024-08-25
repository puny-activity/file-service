package filerepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/queryer"
)

func (r *Repository) GetRootID(ctx context.Context, fileID file.ID) (root.ID, error) {
	return r.getRootID(ctx, r.db, fileID)
}

func (r *Repository) GetRootIDTx(ctx context.Context, tx *sqlx.Tx, fileID file.ID) (root.ID, error) {
	return r.getRootID(ctx, tx, fileID)
}

func (r *Repository) getRootID(ctx context.Context, queryer queryer.Queryer, fileID file.ID) (root.ID, error) {
	query := `
SELECT f.root_id
FROM files f
WHERE f.id = $1
`

	var rootID uuid.UUID

	err := queryer.GetContext(ctx, &rootID, query, fileID.String())
	if err != nil {
		return root.ID{}, err
	}

	return root.ID(rootID), nil
}
