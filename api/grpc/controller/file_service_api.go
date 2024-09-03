package controller

import (
	"github.com/puny-activity/file-service/internal/usecase/fileuc"
	"github.com/puny-activity/file-service/pkg/proto/gen/fileserviceproto"
	"github.com/rs/zerolog"
)

type Controller struct {
	fileserviceproto.FileServiceServer
	fileUseCase *fileuc.UseCase
	log         *zerolog.Logger
}

func New(fileUseCase *fileuc.UseCase, logger *zerolog.Logger) *Controller {
	return &Controller{
		fileUseCase: fileUseCase,
		log:         logger,
	}
}
