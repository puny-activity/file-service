package fileuc

import (
	"context"
	"github.com/golang-module/carbon"
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/filehistory"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
	"io"
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
	ReadFile(ctx context.Context, rootID root.ID, file file.File) (io.ReadCloser, error)
}

type fileRepository interface {
	GetAll(ctx context.Context) ([]file.File, error)
	GetAllByRoot(ctx context.Context, rootID root.ID) ([]file.File, error)
	Create(ctx context.Context, rootID root.ID, fileToCreate file.File) error
	Update(ctx context.Context, fileToUpdate file.File) error
	Delete(ctx context.Context, fileID file.ID) error
	GetRootID(ctx context.Context, fileID file.ID) (root.ID, error)
	Get(ctx context.Context, fileID file.ID) (file.File, error)
	GetAllByIDsTx(ctx context.Context, tx *sqlx.Tx, fileIDs []file.ID) ([]file.File, error)
}

type fileHistoryRepository interface {
	Create(ctx context.Context, row filehistory.Row) error
	GetSinceTx(ctx context.Context, tx *sqlx.Tx, updatedSince carbon.Carbon) ([]filehistory.Row, error)
}
