package rootuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/internal/infrastructure/storage"
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	storageController storageController
	rootRepository    rootRepository
	txManager         txmanager.Transactor
	log               *zerolog.Logger
}

func New(storageController storageController, rootRepository rootRepository, txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		storageController: storageController,
		rootRepository:    rootRepository,
		txManager:         txManager,
		log:               log,
	}
}

type storageController interface {
	Add(cfg storage.Config, rootName string) error
	Remove(rootID root.ID)
	Reset()
	GetRootIDs() []root.ID
}

type rootRepository interface {
	GetAll(ctx context.Context) ([]root.Root, error)
}
