package storage

import (
	"context"
	"errors"
	"github.com/puny-activity/file-service/internal/entity/file"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/pkg/util"
	"github.com/puny-activity/file-service/pkg/werr"
	"github.com/rs/zerolog"
)

type Controller struct {
	storages map[root.ID]Storage
	log      *zerolog.Logger
}

func NewController(log *zerolog.Logger) *Controller {
	return &Controller{
		storages: make(map[root.ID]Storage),
		log:      log,
	}
}

func (c *Controller) Add(cfg Config, log *zerolog.Logger) error {
	newStorage, err := New(cfg, log)
	if err != nil {
		return werr.WrapSE("failed to initialize storage", err)
	}
	c.storages[cfg.ID()] = newStorage
	return nil
}

func (c *Controller) Delete(rootID root.ID) {
	delete(c.storages, rootID)
}

func (c *Controller) Reset() {
	c.storages = make(map[root.ID]Storage)
}

func (c *Controller) GetRootIDs() []root.ID {
	return util.MapKeys(c.storages)
}

func (c *Controller) GetFiles(ctx context.Context, rootID root.ID) ([]file.File, error) {
	storage, ok := c.storages[rootID]
	if !ok {
		return nil, errors.New("storage not found")
	}
	return storage.GetFiles(ctx)
}
