package miniostorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/root"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
)

type Config struct {
	rootID   root.ID
	endpoint string
	username string
	password string
	useSSL   bool
}

type config struct {
	Endpoint *string `json:"endpoint"`
	Username *string `json:"username"`
	Password *string `json:"password"`
	UseSSL   *bool   `json:"use_ssl"`
}

func NewConfig(rootID root.ID, jsonConfig json.RawMessage) (*Config, error) {
	var m config
	err := json.Unmarshal(jsonConfig, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}

	if m.Endpoint == nil {
		return nil, errors.New("endpoint is required")
	}

	if m.Username == nil {
		return nil, errors.New("username is required")
	}

	if m.Password == nil {
		return nil, errors.New("password is required")
	}

	if m.UseSSL == nil {
		return nil, errors.New("use_ssl is required")
	}

	return &Config{
		rootID:   rootID,
		endpoint: *m.Endpoint,
		username: *m.Username,
		password: *m.Password,
		useSSL:   *m.UseSSL,
	}, nil
}

func (c Config) Endpoint() string {
	return c.endpoint
}

func (c Config) Username() string {
	return c.username
}

func (c Config) Password() string {
	return c.password
}

func (c Config) UseSSL() bool {
	return c.useSSL
}

func (c Config) ID() root.ID {
	return c.rootID
}

func (c Config) Type() roottype.Type {
	return roottype.Minio
}

func (c Config) JSONRawMessage() (json.RawMessage, error) {
	return json.Marshal(config{
		Endpoint: &c.endpoint,
		Username: &c.username,
		Password: &c.password,
		UseSSL:   &c.useSSL,
	})
}
