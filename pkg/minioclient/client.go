package minioclient

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Config interface {
	Endpoint() string
	Username() string
	Password() string
	UseSSL() bool
}

func New(cfg Config) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint(), &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Username(), cfg.Password(), ""),
		Secure: cfg.UseSSL(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	return client, nil
}
