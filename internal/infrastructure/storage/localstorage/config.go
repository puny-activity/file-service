package localstorage

import (
	"encoding/json"
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/pkg/werr"
	"os"
)

type Config struct {
	path string
}

type config struct {
	Path string `json:"path"`
}

func NewConfig(jsonConfig json.RawMessage) (*Config, error) {
	var l config
	err := json.Unmarshal(jsonConfig, &l)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}

	_, err = os.Stat(l.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, werr.WrapSE("path does not exist", err)
		}
		return nil, werr.WrapSE("failed to check path", err)
	}

	return &Config{
		path: l.Path,
	}, nil
}

func (c Config) Path() string {
	return c.path
}

func (c Config) Type() roottype.Type {
	return roottype.Local
}

func (c Config) JSONRawMessage() (json.RawMessage, error) {
	return json.Marshal(config{
		Path: c.path,
	})
}
