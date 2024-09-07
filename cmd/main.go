package main

import (
	"context"
	grpcserver "github.com/puny-activity/file-service/api/grpc"
	grpccontroller "github.com/puny-activity/file-service/api/grpc/controller"
	httpcontroller "github.com/puny-activity/file-service/api/http/controller"
	httpmiddleware "github.com/puny-activity/file-service/api/http/middleware"
	httprouter "github.com/puny-activity/file-service/api/http/router"
	"github.com/puny-activity/file-service/config"
	"github.com/puny-activity/file-service/internal/app"
	appconfig "github.com/puny-activity/file-service/internal/config"
	"github.com/puny-activity/file-service/pkg/chimux"
	"github.com/puny-activity/file-service/pkg/httpresp"
	"github.com/puny-activity/file-service/pkg/httpsrvr"
	"github.com/puny-activity/file-service/pkg/werr"
	"github.com/puny-activity/file-service/pkg/zerologger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.Parse()
	if err != nil {
		panic(err)
	}

	log, err := zerologger.NewLogger(cfg.Logger.Level)
	if err != nil {
		panic(err)
	}

	appConfig := appconfig.App{
		Database: appconfig.Database{
			Host:           cfg.App.Database.Host,
			Port:           cfg.App.Database.Port,
			Name:           cfg.App.Database.Name,
			User:           cfg.App.Database.User,
			Password:       cfg.App.Database.Password,
			MigrationsPath: cfg.App.Database.MigrationsPath,
		},
	}

	application := app.New(appConfig, log)
	defer application.Close()

	application.FileUseCase.ScanAll(context.Background())

	chiMux := chimux.New()
	httpMiddleware := httpmiddleware.New()
	httpRespWriter := httpresp.NewWriter()
	httpWrapper := httprouter.NewWrapper(httpRespWriter, nil, log)
	controller := httpcontroller.New(application, httpRespWriter, log)
	httpRouter := httprouter.New(&cfg.API.HTTP, chiMux, httpMiddleware, httpWrapper, controller, log)
	httpRouter.Setup()

	httpServer := httpsrvr.New(
		chiMux,
		httpsrvr.Addr(cfg.API.HTTP.Host, cfg.API.HTTP.Port))

	grpcController := grpccontroller.New(application.FileUseCase, log)
	grpcServer := grpcserver.New(&cfg.API.GRPC, grpcController, log)
	err = grpcServer.Start()
	if err != nil {
		panic(werr.WrapSE("failed to start grpc server", err))
	}
	defer grpcServer.Stop()

	//application.FileUseCase.ScanAll(context.Background())

	//webSocketServer := wsserver.New(&cfg.API.WebSocket, application)
	//webSocketServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info().Str("signal", s.String()).Msg("interrupt")
	}

	err = httpServer.Shutdown()
	if err != nil {
		log.Error().Err(err).Msg("failed to shutdown http server")
	}
	err = application.Close()
	if err != nil {
		log.Error().Err(err).Msg("failed to close application")
	}
}
