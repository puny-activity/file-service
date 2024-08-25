package fileuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/pkg/werr"
	"io"
)

func (u *UseCase) ReadFile(ctx context.Context, fileID file.ID) (io.ReadCloser, error) {
	rootID, err := u.fileRepository.GetRootID(ctx, fileID)
	if err != nil {
		return nil, werr.WrapSE("failed to get root id", err)
	}

	targetFile, err := u.fileRepository.Get(ctx, fileID)
	if err != nil {
		return nil, werr.WrapSE("failed to get file info", err)
	}

	return u.storageController.ReadFile(ctx, rootID, targetFile)
}
