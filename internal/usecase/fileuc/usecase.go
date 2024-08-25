package fileuc

import (
	"context"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	storageController     storageController
	fileRepository        fileRepository
	fileHistoryRepository fileHistoryRepository
	txManager             txmanager.Transactor
	log                   *zerolog.Logger
}

func New(storageController storageController, fileRepository fileRepository,
	fileHistoryRepository fileHistoryRepository, txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		storageController:     storageController,
		fileRepository:        fileRepository,
		fileHistoryRepository: fileHistoryRepository,
		txManager:             txManager,
		log:                   log,
	}
}

type storageController interface {
	GetRootIDs() []root.ID
	GetFiles(ctx context.Context, rootID root.ID) ([]file.File, error)
}

type fileRepository interface {
}

type fileHistoryRepository interface {
}
