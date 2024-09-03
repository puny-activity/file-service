package grpcserver

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/puny-activity/file-service/api/grpc/controller"
	"github.com/puny-activity/file-service/config"
	"github.com/puny-activity/file-service/pkg/proto/gen/fileserviceproto"
	"github.com/puny-activity/file-service/pkg/werr"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct {
	cfg        *config.GRPC
	controller *controller.Controller
	grpcServer *grpc.Server
	log        *zerolog.Logger
}

func New(cfg *config.GRPC, server *controller.Controller, log *zerolog.Logger) *GRPCServer {
	return &GRPCServer{
		cfg:        cfg,
		controller: server,
		log:        log,
	}
}

func (s *GRPCServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", s.cfg.Port))
	if err != nil {
		return werr.WrapSE("failed to listen", err)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor()))
	s.grpcServer = server

	fileserviceproto.RegisterFileServiceServer(server, s.controller)
	go func() {
		err = server.Serve(listener)
		if err != nil {
			panic(werr.WrapSE("failed to serve", err))
		}
	}()

	return nil
}

func (s *GRPCServer) Stop() {
	s.grpcServer.GracefulStop()
}
