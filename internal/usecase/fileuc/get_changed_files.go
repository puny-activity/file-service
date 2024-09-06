package fileuc

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/filehistory/actiontype"
	"github.com/puny-activity/file-service/pkg/werr"
)

func (u *UseCase) GetChangedFiles(ctx context.Context, updatedSince carbon.Carbon) (file.Changes, error) {
	createdFiles := make([]file.File, 0)
	updatedFiles := make([]file.File, 0)
	deletedFileIDs := make([]file.ID, 0)

	err := u.txManager.Transaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		createdFileIDsMap := map[file.ID]struct{}{}
		updatedFileIDsMap := map[file.ID]struct{}{}
		deletedFileIDsMap := map[file.ID]struct{}{}

		historyRows, err := u.fileHistoryRepository.GetSinceTx(ctx, tx, updatedSince)
		if err != nil {
			return werr.WrapSE("failed to get history rows", err)
		}
		for i := range historyRows {
			switch historyRows[i].Action {
			case actiontype.Create:
				createdFileIDsMap[historyRows[i].FileID] = struct{}{}
			case actiontype.Update:
				updatedFileIDsMap[historyRows[i].FileID] = struct{}{}
				_, ok := createdFileIDsMap[historyRows[i].FileID]
				if ok {
					delete(createdFileIDsMap, historyRows[i].FileID)
				}
			case actiontype.Delete:
				deletedFileIDsMap[historyRows[i].FileID] = struct{}{}
				_, ok := createdFileIDsMap[historyRows[i].FileID]
				if ok {
					delete(createdFileIDsMap, historyRows[i].FileID)
				}
				_, ok = updatedFileIDsMap[historyRows[i].FileID]
				if ok {
					delete(updatedFileIDsMap, historyRows[i].FileID)
				}
			default:
				u.log.Warn().Str("actionType", historyRows[i].Action.String()).Msg("unknown action type")
			}
		}

		createdFileIDs := make([]file.ID, 0, len(createdFileIDsMap))
		for fileID := range createdFileIDsMap {
			createdFileIDs = append(createdFileIDs, fileID)
		}
		updatedFileIDs := make([]file.ID, 0, len(updatedFileIDsMap))
		for fileID := range updatedFileIDsMap {
			updatedFileIDs = append(updatedFileIDs, fileID)
		}
		deletedFileIDs = make([]file.ID, 0, len(deletedFileIDsMap))
		for fileID := range deletedFileIDsMap {
			deletedFileIDs = append(deletedFileIDs, fileID)
		}

		createdFiles, err = u.fileRepository.GetAllByIDsTx(ctx, tx, createdFileIDs)
		if err != nil {
			return werr.WrapSE("failed to get all files", err)
		}
		updatedFiles, err = u.fileRepository.GetAllByIDsTx(ctx, tx, updatedFileIDs)
		if err != nil {
			return werr.WrapSE("failed to get all files", err)
		}
		return nil
	})
	if err != nil {
		return file.Changes{}, err
	}

	return file.Changes{
		Created: createdFiles,
		Updated: updatedFiles,
		Deleted: deletedFileIDs,
	}, nil
}
