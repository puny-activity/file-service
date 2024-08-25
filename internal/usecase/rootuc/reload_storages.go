package rootuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/infrastructure/storage"
	"github.com/puny-activity/file-service/pkg/werr"
)

func (u *UseCase) ReloadStorages(ctx context.Context) error {
	u.storageController.Reset()

	roots, err := u.rootRepository.GetAll(ctx)
	if err != nil {
		return werr.WrapSE("failed to get roots", err)
	}

	for i := range roots {
		rootConfig, err := storage.NewConfig(*roots[i].ID, roots[i].Type, roots[i].Config)
		if err != nil {
			return werr.WrapSE("failed to create root config", err)
		}

		err = u.storageController.Add(rootConfig, u.log)
		if err != nil {
			return werr.WrapSE("failed to add storage to controller", err)
		}
	}

	return nil
}
