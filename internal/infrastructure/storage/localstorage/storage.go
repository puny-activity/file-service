package localstorage

import (
	"github.com/rs/zerolog"
)

type Storage struct {
	basePath string
	log      *zerolog.Logger
}

func New(basePath string, log *zerolog.Logger) *Storage {
	return &Storage{
		basePath: basePath,
		log:      log,
	}
}
