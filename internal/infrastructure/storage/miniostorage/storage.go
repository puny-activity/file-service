package miniostorage

import (
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
)

type Storage struct {
	name        string
	minioClient *minio.Client
	log         *zerolog.Logger
}

func New(minioClient *minio.Client, rootName string, log *zerolog.Logger) *Storage {
	return &Storage{
		name:        rootName,
		minioClient: minioClient,
		log:         log,
	}
}

func (s *Storage) Name() string {
	return s.name
}
