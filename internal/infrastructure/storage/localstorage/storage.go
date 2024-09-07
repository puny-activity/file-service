package localstorage

import (
	"github.com/rs/zerolog"
)

type Storage struct {
	name     string
	basePath string
	log      *zerolog.Logger
}

func New(basePath string, rootName string, log *zerolog.Logger) *Storage {
	return &Storage{
		name:     rootName,
		basePath: basePath,
		log:      log,
	}
}

func (s *Storage) Name() string {
	return s.name
}
