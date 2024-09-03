package fileuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/pkg/werr"
)

func (u *UseCase) GetAllFiles(ctx context.Context) ([]file.File, error) {
	files, err := u.fileRepository.GetAll(ctx)
	if err != nil {
		return nil, werr.WrapSE("failed to get all files", err)
	}

	return files, nil
}
