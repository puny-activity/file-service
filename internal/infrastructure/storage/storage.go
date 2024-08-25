package storage

import (
	"context"
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/internal/infrastructure/storage/localstorage"
	"github.com/puny-activity/file-service/internal/infrastructure/storage/miniostorage"
	"github.com/puny-activity/file-service/pkg/minioclient"
	"github.com/puny-activity/file-service/pkg/werr"
	"github.com/rs/zerolog"
	"io"
)

type Storage interface {
	GetFiles(ctx context.Context) ([]file.File, error)
	ReadFile(ctx context.Context, file file.File) (io.ReadCloser, error)
}

func New(cfg Config, log *zerolog.Logger) (Storage, error) {
	switch cfg.Type() {
	case roottype.Local:
		localStorage := localstorage.New(cfg.(*localstorage.Config).Path(), log)
		return localStorage, nil
	case roottype.Minio:
		minioClient, err := minioclient.New(cfg.(*miniostorage.Config))
		if err != nil {
			return nil, werr.WrapSE("failed to create minio client", err)
		}
		minioStorage := miniostorage.New(minioClient, log)
		return minioStorage, nil
	default:
		return nil, fmt.Errorf("unknown root type: %v", cfg.Type())
	}
}
