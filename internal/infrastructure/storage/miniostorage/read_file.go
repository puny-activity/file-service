package miniostorage

import (
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/puny-activity/file-service/internal/entity/file"
	"io"
)

func (s *Storage) ReadFile(ctx context.Context, file file.File) (io.ReadCloser, error) {
	return s.minioClient.GetObject(ctx, bucketPyPath(file.Path), bucketFilenamePyPath(file.Path), minio.GetObjectOptions{})
}
