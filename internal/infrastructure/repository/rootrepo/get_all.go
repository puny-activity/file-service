package rootrepo

import (
	"context"
	"encoding/json"
	"github.com/golang-module/carbon"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/pkg/queryer"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/puny-activity/file-service/pkg/werr"
)

type getAllEntity struct {
	ID        uuid.UUID       `db:"id"`
	Type      string          `db:"type"`
	Name      string          `db:"name"`
	Config    json.RawMessage `db:"config"`
	ScannedAt *string         `db:"scanned_at"`
}

func (r *Repository) GetAll(ctx context.Context) ([]root.Root, error) {
	return r.getAll(ctx, r.db)
}

func (r *Repository) GetAllTx(ctx context.Context, tx *sqlx.Tx) ([]root.Root, error) {
	return r.getAll(ctx, tx)
}

func (r *Repository) getAll(ctx context.Context, queryer queryer.Queryer) ([]root.Root, error) {
	query := `
SELECT r.id,
       r.type,
       r.name,
       r.config,
       r.scanned_at
FROM roots r
`

	rootsRepo := make([]getAllEntity, 0)
	err := queryer.SelectContext(ctx, &rootsRepo, query)
	if err != nil {
		return nil, err
	}

	roots := make([]root.Root, len(rootsRepo))
	for i := range rootsRepo {
		rootType, err := roottype.Parse(rootsRepo[i].Type)
		if err != nil {
			return nil, werr.WrapSE("failed to parse root type", err)
		}

		var scannedAt *carbon.Carbon = nil
		if rootsRepo[i].ScannedAt != nil {
			scannedAt = util.ToPointer(carbon.Parse(*rootsRepo[i].ScannedAt))
		}

		roots[i] = root.Root{
			ID:        util.ToPointer(root.ID(rootsRepo[i].ID)),
			Type:      rootType,
			Name:      rootsRepo[i].Name,
			Config:    rootsRepo[i].Config,
			ScannedAt: scannedAt,
		}
	}

	return roots, nil
}
