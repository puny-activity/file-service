package app

import (
	"context"
	"encoding/json"
	"github.com/puny-activity/file-service/internal/config"
	"github.com/puny-activity/file-service/internal/entity/root/roottype"
	"github.com/puny-activity/file-service/internal/infrastructure/storage"
	"github.com/puny-activity/file-service/internal/usecase/rootuc"
	"github.com/puny-activity/file-service/pkg/postgres"
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/rs/zerolog"
)

type App struct {
	RootUseCase *rootuc.UseCase
	log         *zerolog.Logger
}

func New(cfg config.App, log *zerolog.Logger) *App {
	db, err := postgres.New(cfg.Database.ConnectionString())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.RunMigrations(cfg.Database.MigrationsPath)
	if err != nil {
		panic(err)
	}

	txManager := txmanager.New(db.DB)

	localStorageConfig, err := storage.NewConfig(roottype.Local, json.RawMessage(
		"{ \"path\": \"/home/zalimannard/Музыка/\" }",
	))
	if err != nil {
		panic(err)
	}
	localStorage, err := storage.New(localStorageConfig, log)
	if err != nil {
		panic(err)
	}

	minioStorageConfig, err := storage.NewConfig(roottype.Minio, json.RawMessage(
		"{ \"endpoint\": \"localhost:9000\", \"username\": \"user\", \"password\": \"password\", \"use_ssl\": false }",
	))
	if err != nil {
		panic(err)
	}
	minioStorage, err := storage.New(minioStorageConfig, log)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	localFiles, err := localStorage.GetFiles(ctx)
	if err != nil {
		panic(err)
	}
	log.Info().Int("count", len(localFiles)).Any("localFiles", localFiles).Msg("local files")

	minioFiles, err := minioStorage.GetFiles(ctx)
	if err != nil {
		panic(err)
	}
	log.Info().Int("count", len(minioFiles)).Any("minioFiles", minioFiles).Msg("minio files")

	rootUseCase := rootuc.New(txManager, log)

	return &App{
		RootUseCase: rootUseCase,
		log:         log,
	}
}

func (a *App) Close() error {
	return nil
}
