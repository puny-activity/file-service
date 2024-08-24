package rootrepo

import (
	"github.com/jmoiron/sqlx"
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
)

type Repository struct {
	db        *sqlx.Tx
	txManager txmanager.Transactor
	log       *zerolog.Logger
}
