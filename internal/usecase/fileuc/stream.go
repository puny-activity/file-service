package fileuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/pkg/werr"
	"io"
)

func (u *UseCase) StreamFile(ctx context.Context, fileID file.ID) (io.ReadCloser, file.File, error) {
	rootID, err := u.fileRepository.GetRootID(ctx, fileID)
	if err != nil {
		return nil, file.File{}, werr.WrapSE("failed to get root id", err)
	}

	targetFile, err := u.fileRepository.Get(ctx, fileID)
	if err != nil {
		return nil, file.File{}, werr.WrapSE("failed to get file info", err)
	}

	fileReader, err := u.storageController.ReadFile(ctx, rootID, targetFile)
	if err != nil {
		return nil, file.File{}, werr.WrapSE("failed to read file", err)
	}

	return fileReader, targetFile, nil
}
