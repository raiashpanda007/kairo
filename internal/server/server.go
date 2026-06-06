package server

import (
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	server *grpc.Server
	logger *zap.Logger
	addr   string
}

func New(addr string, logger *zap.Logger) *GRPCServer {
	srv := grpc.NewServer()
	reflection.Register(srv)

	return &GRPCServer{
		server: srv,
		logger: logger,
		addr:   addr,
	}
}

func (s *GRPCServer) Server() *grpc.Server {
	return s.server
}

func (s *GRPCServer) Run() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.logger.Info("gRPC server listening", zap.String("addr", s.addr))
	return s.server.Serve(lis)
}

func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}
