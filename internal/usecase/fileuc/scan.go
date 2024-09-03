package fileuc

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/filehistory"
	"github.com/puny-activity/file-service/internal/entity/filehistory/actiontype"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/werr"
)

func (u *UseCase) Scan(ctx context.Context, rootID root.ID) {
	err := u.scan(ctx, rootID)
	if err != nil {
		u.log.Warn().Err(err).Str("rootId", rootID.String()).Msg("failed to scan root")
	}
}

func (u *UseCase) ScanAll(ctx context.Context) {
	rootIDs := u.storageController.GetRootIDs()
	for i := range rootIDs {
		err := u.scan(ctx, rootIDs[i])
		if err != nil {
			u.log.Warn().Err(err).Str("rootId", rootIDs[i].String()).Msg("failed to scan root")
		}
	}
}

func (u *UseCase) scan(ctx context.Context, rootID root.ID) error {
	u.log.Info().Str("rootId", rootID.String()).Msg("scanning root")

	savedFilesList, err := u.fileRepository.GetAllByRoot(ctx, rootID)
	if err != nil {
		return werr.WrapSE("failed to get saved files", err)
	}
	savedFilesByPath := make(map[string]file.File)
	for i := range savedFilesList {
		savedFilesByPath[savedFilesList[i].Path] = savedFilesList[i]
	}

	actualFilesList, err := u.storageController.GetFiles(ctx, rootID)
	if err != nil {
		return werr.WrapSE("failed to get actual files", err)
	}
	actualFilesByPath := make(map[string]file.File)
	for i := range actualFilesList {
		actualFilesByPath[actualFilesList[i].Path] = actualFilesList[i]
	}

	filesToCreate := make([]file.File, 0)
	filesToUpdate := make([]file.File, 0)
	filesToDelete := make([]file.File, 0)

	for savedFilePath, savedFile := range savedFilesByPath {
		actualFile, ok := actualFilesByPath[savedFilePath]
		if !ok {
			filesToDelete = append(filesToDelete, savedFile)
			continue
		}
		if savedFile.MD5 != actualFile.MD5 {
			filesToUpdate = append(filesToUpdate, savedFile)
		}
	}

	for actualFilePath, actualFile := range actualFilesByPath {
		_, ok := savedFilesByPath[actualFilePath]
		if !ok {
			filesToCreate = append(filesToCreate, actualFile)
		}
	}

	for _, fileToCreate := range filesToCreate {
		fileToCreate = fileToCreate.GenerateID()
		err := u.fileRepository.Create(ctx, rootID, fileToCreate)
		if err != nil {
			u.log.Warn().Str("path", fileToCreate.Path).Err(err).Msg("failed to create file")
			continue
		}

		historyRow := filehistory.Row{
			FileID:      *fileToCreate.ID,
			Action:      actiontype.Create,
			PerformedAt: carbon.Now(),
		}
		historyRow = historyRow.GenerateID()
		err = u.fileHistoryRepository.Create(ctx, historyRow)
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to create file history row")
		}
	}

	for _, fileToUpdate := range filesToUpdate {
		err := u.fileRepository.Update(ctx, fileToUpdate)
		if err != nil {
			u.log.Warn().Err(err).Str("path", fileToUpdate.Path).Msg("failed to update file")
			continue
		}

		historyRow := filehistory.Row{
			FileID:      *fileToUpdate.ID,
			Action:      actiontype.Update,
			PerformedAt: carbon.Now(),
		}
		historyRow = historyRow.GenerateID()
		err = u.fileHistoryRepository.Create(ctx, historyRow)
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to create file history row")
		}
	}

	for _, fileToDelete := range filesToDelete {
		err := u.fileRepository.Delete(ctx, *fileToDelete.ID)
		if err != nil {
			u.log.Warn().Err(err).Str("path", fileToDelete.Path).Msg("failed to delete file")
			continue
		}

		historyRow := filehistory.Row{
			FileID:      *fileToDelete.ID,
			Action:      actiontype.Delete,
			PerformedAt: carbon.Now(),
		}
		historyRow = historyRow.GenerateID()
		err = u.fileHistoryRepository.Create(ctx, historyRow)
		if err != nil {
			u.log.Warn().Err(err).Msg("failed to create file history row")
		}
	}

	return nil
}
