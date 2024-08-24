package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/internal/infrastructure/storage/localstorage"
	"github.com/puny-activity/file-service/internal/infrastructure/storage/miniostorage"
	"github.com/puny-activity/file-service/pkg/minioclient"
	"github.com/puny-activity/file-service/pkg/werr"
	"github.com/rs/zerolog"
)

type Config interface {
	Type() roottype.Type
	JSONRawMessage() (json.RawMessage, error)
}

func NewConfig(rootType roottype.Type, jsonConfig json.RawMessage) (Config, error) {
	switch rootType {
	case roottype.Local:
		config, err := localstorage.NewConfig(jsonConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to initialize local storage", err)
		}
		return config, nil
	case roottype.Minio:
		config, err := miniostorage.NewConfig(jsonConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to initialize minio storage", err)
		}
		return config, nil
	default:
		return nil, errors.New("unknown root type")
	}
}

type Storage interface {
	GetFiles(ctx context.Context) ([]file.File, error)
}

func New(cfg Config, log *zerolog.Logger) (Storage, error) {
	jsonConfig, err := cfg.JSONRawMessage()
	if err != nil {
		return nil, werr.WrapSE("failed to load json config", err)
	}

	switch cfg.Type() {
	case roottype.Local:
		localConfig, err := localstorage.NewConfig(jsonConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to create local config", err)
		}
		localStorage := localstorage.New(localConfig.Path(), log)
		return localStorage, nil
	case roottype.Minio:
		minioConfig, err := miniostorage.NewConfig(jsonConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to create minio config", err)
		}
		minioClient, err := minioclient.New(minioConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to create minio client", err)
		}
		minioStorage := miniostorage.New(minioClient, log)
		return minioStorage, nil
	default:
		return nil, fmt.Errorf("unknown root type: %v", cfg.Type())
	}
}
