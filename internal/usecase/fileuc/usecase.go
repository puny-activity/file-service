package fileuc

import (
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	fileRepository        fileRepository
	fileHistoryRepository fileHistoryRepository
	txManager             txmanager.Transactor
	log                   *zerolog.Logger
}

func New(txManager txmanager.Transactor, fileRepository fileRepository, fileHistoryRepository fileHistoryRepository,
	log *zerolog.Logger) *UseCase {
	return &UseCase{
		txManager: txManager,
		log:       log,
	}
}

type fileRepository interface {
}

type fileHistoryRepository interface {
}
