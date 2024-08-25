package filehistoryrepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/filehistory"
	"github.com/puny-activity/file-service/pkg/queryer"
)

type createEntity struct {
	ID          uuid.UUID `db:"id"`
	FileID      uuid.UUID `db:"file_id"`
	ActionType  string    `db:"action_type"`
	PerformedAt string    `db:"performed_at"`
}

func (r *Repository) Create(ctx context.Context, row filehistory.Row) error {
	return r.create(ctx, r.db, row)
}

func (r *Repository) CreateTx(ctx context.Context, tx *sqlx.Tx, row filehistory.Row) error {
	return r.create(ctx, tx, row)
}

func (r *Repository) create(ctx context.Context, queryer queryer.Queryer, row filehistory.Row) error {
	query := `
INSERT INTO file_history (id, file_id, action_type, performed_at)
VALUES (:id, :file_id, :action_type, :performed_at)
`

	parameter := createEntity{
		ID:          uuid.UUID(*row.ID),
		FileID:      uuid.UUID(row.FileID),
		ActionType:  row.Action.String(),
		PerformedAt: row.PerformedAt.String(),
	}

	_, err := queryer.NamedExecContext(ctx, query, parameter)
	if err != nil {
		return err
	}

	return nil
}
