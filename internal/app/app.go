package app

import (
	"context"
	"github.com/puny-activity/file-service/internal/config"
	"github.com/puny-activity/file-service/internal/infrastructure/repository/filehistoryrepo"
	"github.com/puny-activity/file-service/internal/infrastructure/repository/filerepo"
	"github.com/puny-activity/file-service/internal/infrastructure/repository/rootrepo"
	"github.com/puny-activity/file-service/internal/infrastructure/storage"
	"github.com/puny-activity/file-service/internal/usecase/fileuc"
	"github.com/puny-activity/file-service/internal/usecase/rootuc"
	"github.com/puny-activity/file-service/pkg/postgres"
	"github.com/puny-activity/file-service/pkg/txmanager"
	"github.com/puny-activity/file-service/pkg/werr"
	"github.com/rs/zerolog"
)

type App struct {
	RootUseCase *rootuc.UseCase
	FileUseCase *fileuc.UseCase
	db          *postgres.Postgres
	log         *zerolog.Logger
}

func New(cfg config.App, log *zerolog.Logger) *App {
	db, err := postgres.New(cfg.Database.ConnectionString())
	if err != nil {
		panic(err)
	}
	err = db.RunMigrations(cfg.Database.MigrationsPath)
	if err != nil {
		panic(err)
	}

	txManager := txmanager.New(db.DB)

	storageController := storage.NewController(log)

	rootRepository := rootrepo.New(db.DB, txManager, log)
	fileRepository := filerepo.New(db.DB, txManager, log)
	fileHistoryRepository := filehistoryrepo.New(db.DB, txManager, log)

	rootUseCase := rootuc.New(storageController, rootRepository, txManager, log)
	fileUseCase := fileuc.New(storageController, fileRepository, fileHistoryRepository, txManager, log)

	err = rootUseCase.ReloadStorages(context.Background())
	if err != nil {
		panic(err)
	}

	return &App{
		RootUseCase: rootUseCase,
		FileUseCase: fileUseCase,
		db:          db,
		log:         log,
	}
}

func (a *App) Close() error {
	err := a.db.Close()
	if err != nil {
		return werr.WrapSE("failed to close database connection", err)
	}
	return nil
}
