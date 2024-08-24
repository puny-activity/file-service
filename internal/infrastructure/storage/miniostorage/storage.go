package miniostorage

import (
	"github.com/minio/minio-go/v7"
	"github.com/rs/zerolog"
)

type Storage struct {
	minioClient *minio.Client
	log         *zerolog.Logger
}

func New(minioClient *minio.Client, log *zerolog.Logger) *Storage {
	return &Storage{
		minioClient: minioClient,
		log:         log,
	}
}
