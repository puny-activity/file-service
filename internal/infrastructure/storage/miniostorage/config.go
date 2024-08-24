package miniostorage

import (
	"encoding/json"
	"fmt"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
)

type Config struct {
	endpoint string
	username string
	password string
	useSSL   bool
}

type config struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
	UseSSL   bool   `json:"use_ssl"`
}

func NewConfig(jsonConfig json.RawMessage) (*Config, error) {
	var m config
	err := json.Unmarshal(jsonConfig, &m)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json: %v", err)
	}
	return &Config{
		endpoint: m.Endpoint,
		username: m.Username,
		password: m.Password,
		useSSL:   m.UseSSL,
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

func (c Config) Type() roottype.Type {
	return roottype.Minio
}

func (c Config) JSONRawMessage() (json.RawMessage, error) {
	return json.Marshal(config{
		Endpoint: c.endpoint,
		Username: c.username,
		Password: c.password,
		UseSSL:   c.useSSL,
	})
}
