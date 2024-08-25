package storage

import (
	"encoding/json"
	"errors"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/internal/infrastructure/storage/localstorage"
	"github.com/puny-activity/file-service/internal/infrastructure/storage/miniostorage"
	"github.com/puny-activity/file-service/pkg/werr"
)

type Config interface {
	ID() root.ID
	Type() roottype.Type
	JSONRawMessage() (json.RawMessage, error)
}

func NewConfig(rootID root.ID, rootType roottype.Type, jsonConfig json.RawMessage) (Config, error) {
	switch rootType {
	case roottype.Local:
		config, err := localstorage.NewConfig(rootID, jsonConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to initialize local storage", err)
		}
		return config, nil
	case roottype.Minio:
		config, err := miniostorage.NewConfig(rootID, jsonConfig)
		if err != nil {
			return nil, werr.WrapSE("failed to initialize minio storage", err)
		}
		return config, nil
	default:
		return nil, errors.New("unknown root type")
	}
}
