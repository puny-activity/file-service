package rootuc

import (
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
)

type UseCase struct {
	txManager txmanager.Transactor
	log       *zerolog.Logger
}

func New(txManager txmanager.Transactor, log *zerolog.Logger) *UseCase {
	return &UseCase{
		txManager: txManager,
		log:       log,
	}
}
