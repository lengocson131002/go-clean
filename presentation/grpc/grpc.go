package gprc

import (
	"context"
	"fmt"
	"net"

	"github.com/lengocson131002/go-clean/bootstrap"
	"github.com/lengocson131002/go-clean/pkg/logger"
	pb "github.com/lengocson131002/go-clean/presentation/grpc/pb"
	"github.com/lengocson131002/go-clean/presentation/grpc/server"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	cfg    *bootstrap.ServerConfig
	logger logger.Logger
}

func NewGrpcServer(cfg *bootstrap.ServerConfig, logger logger.Logger) *GrpcServer {
	return &GrpcServer{
		cfg:    cfg,
		logger: logger,
	}
}

func (s *GrpcServer) Start(ctx context.Context) error {
	network, gPort := "tcp", s.cfg.GrpcPort
	lis, err := net.Listen(network, fmt.Sprintf("localhost:%d", gPort))

	if err != nil {
		return err
	}

	gSrv := grpc.NewServer()
	tSrv := server.NewT24AccountServer(s.logger)
	pb.RegisterT24AccountServiceServer(gSrv, tSrv)

	if err != nil {
		return err
	}

	go func() {
		defer func() {
			gSrv.GracefulStop()
			s.logger.Info(ctx, "Stop GRPC Server")
		}()
		<-ctx.Done()
	}()

	s.logger.Infof(ctx, "Start GRPC server at port: %v", gPort)
	if err := gSrv.Serve(lis); err != nil {
		return fmt.Errorf("Failed to serve GRPC %w", err)
	}

	return nil
}
